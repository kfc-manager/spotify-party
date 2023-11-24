from selenium import webdriver
from tempfile import mkdtemp
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions
import urllib.parse
import boto3
import os
import json
import time
import random


def handler(event=None, context=None):

    # create AWS SecretsManager client
    client = boto3.client("secretsmanager")

    # retrieve secrets from SecretsManager
    response = client.get_secret_value(SecretId=os.environ["STATIC_SECRETS_ID"])
    staticSecrets = json.loads(response["SecretString"])

    # build URL to for Spotify API authorization
    params = {
        "response_type": "code",
        "client_id": staticSecrets["client_id"],
        "scope": "user-read-currently-playing user-read-playback-state user-modify-playback-state",
        "redirect_uri": os.environ["REDIRECT_URI"],
    }
    url = "https://accounts.spotify.com/authorize?" + urllib.parse.urlencode(params)

    # initialize Chrome Webdriver
    options = webdriver.ChromeOptions()
    service = webdriver.ChromeService("/opt/chromedriver")
    options.binary_location = "/opt/chrome/chrome"
    options.add_argument("--headless=new")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-gpu")
    options.add_argument("--window-size=1280x1696")
    options.add_argument("--single-process")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-dev-tools")
    options.add_argument("--no-zygote")
    options.add_argument(f"--user-data-dir={mkdtemp()}")
    options.add_argument(f"--data-path={mkdtemp()}")
    options.add_argument(f"--disk-cache-dir={mkdtemp()}")
    options.add_argument("--remote-debugging-port=9222")

    # to avoid captchas
    options.add_argument("--disable-blink-features=AutomationControlled")
    options.add_experimental_option("excludeSwitches", ["enable-automation"])
    options.add_experimental_option('useAutomationExtension', False)

    chrome = webdriver.Chrome(options=options, service=service)

    # to avoid captchas
    chrome.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")

    # visit Spotify API URL and sign in with Spotify credentials
    chrome.get(url)
    try:
        time.sleep(round(random.uniform(0.5, 3)))
        usernameInput = chrome.find_element(By.XPATH, "//input[@placeholder='Email or username']")
        usernameInput.send_keys(staticSecrets["spotify_username"])
        time.sleep(round(random.uniform(0.5, 3)))
        passwordInput = chrome.find_element(By.XPATH, "//input[@placeholder='Password']")
        passwordInput.send_keys(staticSecrets["spotify_password"])
        time.sleep(round(random.uniform(0.5, 3)))
        signInButton = chrome.find_element(By.XPATH, "//*[contains(text(), 'Log In')]")
        signInButton.click()
    except Exception as e:
        raise e

    # wait until the redirect happened
    try:
        wait = WebDriverWait(chrome, 20)
        wait.until(expected_conditions.url_contains(os.environ["REDIRECT_URI"]))
    except Exception:
        raise Exception(chrome.find_element(By.XPATH, value="//html").text)

    # check if the redirect to callback enpoint was successful
    result = chrome.find_element(By.XPATH, value="//html").text
    if result != "Success":
        raise Exception("Redirect was not successful and was not able to retrieve access token")

    return {"statusCode": 200, "body": "Success"}
