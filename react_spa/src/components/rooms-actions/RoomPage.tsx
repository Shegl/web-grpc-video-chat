import {createContext, useContext, useEffect, useState} from 'react';
import {Fallback, Logout, useAuth, UserContext} from "../../App";
import { Button, Col, Container, Form, Row } from "react-bootstrap";
import { useCookies } from "react-cookie";
import { useNavigate } from 'react-router-dom';
import axios from "axios";

const roomCheck = (cookies, setCookie, navigate, context, setLoaded) => {
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
                setAuthenticated(false);
                setCookie('userUuid', '', { path: '/' });
                navigate('/');
            })
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
        axios.post('http://dev.test:3000/room/leave', {uuid: context.userData.uuid}).then(
            navigate('/home')
        ).catch(
            navigate('/home')
        )
    };

    return (
        !loaded ? <Fallback/> :
        <>
            <Container className="p-1 mb-2 bg-light rounded-3">
                <h4 className="Header">Room uuid: {context.userData.roomUuid}</h4>
                <Button type="button" variant="warning" onClick={handleLeave} className="btn-lg">Leave&nbsp;room</Button>
            </Container>
            <Form>
                <div className="card">
                    <Container>
                        <Row>
                            <Col xs={3} className="text-center">
                                <Button type="submit" variant="success" className="btn-lg">Create&nbsp;room</Button>
                            </Col>
                            <Col xs={1} className="text-center">
                                <p className="some-pad-top">or</p>
                            </Col>
                            <Col xs={7} className="text-center">
                                <Form.Control className="some-margin-top" type="text" name="roomId" id="roomId"></Form.Control>
                            </Col>
                            <Col xs={1}>
                                <Button className="some-margin-top" type="submit" variant="primary">Join</Button>
                            </Col>
                        </Row>
                    </Container>
                </div>
            </Form>
            <p className="read-the-docs">
                Create room or join existing one
            </p>
        </>
    )
}

export default RoomPage