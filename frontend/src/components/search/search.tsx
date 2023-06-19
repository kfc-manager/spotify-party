import { MouseEventHandler, useEffect, useState } from "react"
import "./search.css"
import { ApiResponse, ApiSong } from "../../types"
import Song from "../song/song"
import {MdOutlineAdd} from "react-icons/md"

const Search = (): JSX.Element => {

    const [query, setQuery] = useState<string>("")
    const [result, setResult] = useState<ApiSong[]>([])
    const [message, setMessage] = useState<string>("")

    const handleChange = (event: any): void => {
        setQuery(event.target.value)
    }

    // factoring for each button a custom onClick function
    // with the id of the song in its row
    const handleClick = (id: string): MouseEventHandler<HTMLButtonElement> => {
        return (): void => {
            if (id.length < 1) {
                setMessage("The song id was corrupted")
                setQuery("")
                setResult([])
                return 
            }
            fetch("https://spotify-party-zty7jo4vkq-ey.a.run.app/queue/" + id, { method: "POST" })
                .then((res: Response): Promise<ApiResponse> => {
                    if (res.status === 204) {
                        setMessage("Song has been added to the queue.")
                        setQuery("")
                        setResult([])
                        return new Promise((): ApiResponse => {
                            return {} as ApiResponse
                        })
                    }
                    return res.json()
                }).then((res: ApiResponse): void => {
                    if (res.error !== undefined) {
                        setMessage(res.error.status + " " + res.error.message)
                        setQuery("")
                        setResult([])
                    }
                })
        }
    }

    useEffect((): void => {
        if (query.length < 1) {
            return
        }
        fetch("https://spotify-party-zty7jo4vkq-ey.a.run.app/search/" + query).then((res: Response): Promise<ApiResponse> => {
            return res.json()
        }).then((res: ApiResponse): void => {
            if (res.error !== undefined) {
                setMessage(res.error.status + " " + res.error.message)
                setResult([])
                return
            }
            if (res.result.items === null) {
                setMessage("No songs have been found.")
                setResult([])
                return
            }
            setMessage("")
            setResult(res.result.items)
        })
    }, [query])

    return <div className="search">
        <input type="text" placeholder="Search.." onChange={handleChange} className="search-bar"/>
        <div className="search-item-container">
            {result.map((song: ApiSong): JSX.Element => {
                return <div className="item-row-container" key={result.indexOf(song)}>
                    <Song item={song} search={false}/>
                    <div className="add-button-container">
                        <button className="add-button" onClick={handleClick(song.id)}><MdOutlineAdd className="add-button-icon"/></button>
                    </div>
                </div>
            })}
        </div>
        <div className="search-message">{message}</div>
    </div>
}

export default Search
