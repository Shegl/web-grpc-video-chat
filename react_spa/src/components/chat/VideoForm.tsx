import { useContext, useEffect, useRef } from 'react';
import { Button, Card, InputGroup} from "react-bootstrap";
import { UserContext } from "../../App";

const VideoForm = (props: {owner: boolean;}) => {
    let owner = props.owner
    const { userData } = useContext(UserContext);
    const videoElement = useRef<HTMLVideoElement>(null);

    const captureAndSend = async (stream) => {
        let videoTrack = stream.getVideoTracks()[0];
        const blob = await videoTrack.captureFrame();
    }

    useEffect(() => {
        if (owner) {
            navigator.mediaDevices.getUserMedia({ video: true, audio: true })
                .then(stream => {
                    if (videoElement.current) {
                        videoElement.current.srcObject = stream;
                        videoElement.current.play().then(() => {
                            // lets send stream to server

                        });
                    }
                })
                .catch(err => {
                    // alert(`Following error occured: ${err}`);
                });
        }
    },[userData]);

    return (
        <>
            <Card className="video-card-canvas">
                { owner ?
                    <Card.Header key="you">You ({userData.username})</Card.Header>
                    : <Card.Header key="guest">Guest</Card.Header>
                }
                <video ref={videoElement} className="video-canvas" autoPlay playsInline></video>
                { owner ? <Card.Footer>
                    <InputGroup>
                        <Button variant="secondary">Mute</Button>
                        <Button variant="secondary">Cam turn off</Button>
                        <Button variant="secondary">Blur background</Button>
                    </InputGroup>
                </Card.Footer> : <Card.Footer><Button variant="secondary">Kick</Button></Card.Footer> }
            </Card>
        </>
    );
};

export default VideoForm;