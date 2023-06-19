import { useEffect, useState } from "react"
import "./queue.css"
import { ApiResponse, ApiSong } from "../../types"
import Song from "../song/song"


const Queue = (): JSX.Element => {

    const [queue, setQueue] = useState<ApiSong[]>([])
    const [error, setError] = useState<string>("")

    useEffect(() => {
        fetch("https://spotify-party-zty7jo4vkq-ey.a.run.app/queue").then((res: Response): Promise<ApiResponse> => {
            return res.json()
        }).then((res: ApiResponse): void => {
            if (res.error !== undefined) {
                setError(res.error.status + " " + res.error.message)
                setQueue([])
                return
            }
            if (res.queue.items === null) {
                // user is currently not listening to music
                setError("User is currently not listening to music.")
                setQueue([])
                return
            }
            setError("")
            setQueue(res.queue.items)
        })
    }, [])

    return <div className="queue">
        <div className="queue-item-container">
            {queue.map((song: ApiSong): JSX.Element => {
                return <Song key={queue.indexOf(song)} item={song} search={false}/>
            })}
            <div className="queue-error">{error}</div>
        </div>
    </div>
}

export default Queue


