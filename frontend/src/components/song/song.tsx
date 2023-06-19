import { ApiResponse, ApiSong } from "../../types"
import "./song.css"
import {MdOutlineAdd} from "react-icons/md"

const Song = (props: { item: ApiSong, search: boolean }): JSX.Element => {

    const buildArtistString = (artists: string[]): string => {

        let str: string = artists[0]
        artists.slice(1,artists.length).map((artist: string): void => {
            str += ", " + artist
        })
        
        return str
    }

    const toMinutes = (duration: number): string => {

        let minutes: number = Math.floor(duration/60000)
        let seconds: string = (duration%60000) + "";

        while (seconds.length < 5) {
            seconds = "0" + seconds;
        }

        return minutes + ":" + seconds.slice(0,2)
    }

    const breakString = (name: string): string => {

        if (name.length > 40) {
            return name.slice(0, 39) + "..."
        }

        return name
    }

    const renderButton = (): JSX.Element => {
        if (props.search) {
            return <button className="add-button2" onClick={handleClick}><MdOutlineAdd className="add-button-icon2"/></button>
        }
        return <div/>
    }

    const handleClick = (): void => {
        if (props.item.id.length < 1) {
            // TODO error
            return 
        }
        fetch("https://spotify-party-zty7jo4vkq-ey.a.run.app/queue/" + props.item.id, { method: "POST" })
            .then((res: Response): Promise<ApiResponse> => {
                if (res.status === 204) {
                    window.location.replace("https://spotify-party-zty7jo4vkq-ey.a.run.app") // i was lazy (button click should be signaled to user)
                }
                return res.json()
            }).then((res: ApiResponse): void => {
                if (res.error !== undefined) {
                    // TODO
                }
            })
    }

    return <div className="song">
        <img className="icon" src={ props.item.image_url } alt=""></img>
            <div className="information-container">
                <div className="name">{ breakString(props.item.name) }</div>
                <div className="artists">{ breakString(buildArtistString(props.item.artists)) }</div>
            </div>
            <div className="song-end-container">
                {renderButton()}
                <div className="duration">{ toMinutes(props.item.duration_ms) }</div>
            </div>
    </div>
}

export default Song
