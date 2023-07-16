import { useContext, useEffect, useState } from 'react';
import { Fallback, UserContext } from "../../App";
import {Button, Col, Container, Form, Row, InputGroup} from "react-bootstrap";
import { useCookies } from "react-cookie";
import { useNavigate } from 'react-router-dom';
import axios from "axios";
import ChatForm from "../chat/ChatForm";
import VideoForm from "../chat/VideoForm";

const roomCheck = (cookies: any, setCookie: any, navigate: any, context: any, setLoaded: any) => {
    const { authenticated, setAuthenticated, userData, setUserData } = context;
    if (authenticated && userData.inRoom) {
        setLoaded(true);
    } else {
        let userUUID = cookies.userUuid;
        if (userUUID) {
            axios.post('http://dev.test:3000/room/state', {uuid: userUUID}).then(
                (response) => {
                    if (response.status == 200) {
                        let userData = context.userData;
                        userData.uuid = userUUID;
                        userData.inRoom = true;
                        userData.roomAuthor = response.data.author.uuid == context.userData.uuid;
                        userData.username = userData.roomAuthor ? response.data.author.username : response.data.guest.username;
                        userData.roomUuid = response.data.uuid
                        setAuthenticated(true);
                        setUserData(userData);
                        setLoaded(true);
                        return;
                    } else {
                        navigate('/home');
                    }
                }
            ).catch((error) => {
                if (error.response) {
                    if (error.response.status == 422) {
                        navigate('/home');
                        return;
                    }
                }
                setCookie('userUuid', '', { path: '/' });
                navigate('/');
            })
        } else {
            navigate('/');
        }
    }
}

function RoomPage() {
    const [loaded, setLoaded] = useState(false);
    const context = useContext(UserContext);
    const [cookies, setCookie] = useCookies(['userUuid']);
    const navigate = useNavigate();

    useEffect(() => {
        roomCheck(cookies, setCookie, navigate, context, setLoaded);
    });

    const handleLeave = () => {
        axios.post('http://dev.test:3000/room/leave', {uuid: context.userData.uuid}).then((_) => {
            navigate('/home');
        }).catch((_) => {
            navigate('/home');
        });
    };

    return (
        !loaded ? <Fallback/> :
        <>
            <Container>
                <Row>
                    <InputGroup className="some-margin-bottom">
                        <InputGroup.Text>
                            Room uuid:
                        </InputGroup.Text>
                        <Button variant="secondary">Copy</Button>
                        <Form.Control
                            type="text"
                            value={context.userData.roomUuid}
                            placeholder={context.userData.roomUuid}
                            aria-label={context.userData.roomUuid}
                            aria-describedby="btnGroupAddon" readOnly={true}
                        />
                        <Button type="button" variant="warning" onClick={handleLeave} >Leave&nbsp;room</Button>
                    </InputGroup>
                </Row>
                <Row>
                    <Col className="text-center">
                        <VideoForm owner={context.userData.roomAuthor}/>
                    </Col>
                    <Col className="text-center">
                        <VideoForm owner={!context.userData.roomAuthor}/>
                    </Col>
                    <Col className="text-center">
                        <ChatForm/>
                    </Col>
                </Row>
            </Container>
        </>
    )
}

export default RoomPage