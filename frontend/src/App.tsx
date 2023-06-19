import "./App.css"
import { useState } from "react"
import Queue from "./components/queue/queue"
import Search from "./components/search/search"

const App = (): JSX.Element => {

    const [mode, setMode] = useState<boolean>(true)

    const renderComponent = (): JSX.Element => {
        if (mode) {
            return <Queue/>
        }
        return <Search/>
    }

    return (
        <div className="app">
            <div className="button-container">
                <button className="app-button" onClick={(): void => {setMode(true)}}>Queue</button>
                <button className="app-button" onClick={(): void => {setMode(false)}}>Search</button>
            </div>
            {renderComponent()}
        </div>
    )
}

export default App
