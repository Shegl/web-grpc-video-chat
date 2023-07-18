import { useContext, useEffect, useState } from 'react';
import { Fallback, Logout, useAuth, UserContext } from "../../App";
import { Button, Col, Container, Form, Row } from "react-bootstrap";
import { useCookies } from "react-cookie";
import { useNavigate } from 'react-router-dom';
import axios from "axios";

function HomePage() {
    const [loaded, setLoaded] = useState(false);
    const context = useContext(UserContext);
    const [cookies, setCookie] = useCookies(['userUuid']);
    const navigate = useNavigate();

    useEffect(() => {
        useAuth("/home", false, cookies, setCookie, navigate, context, setLoaded);
    });

    const handleClickCreateRoom = async () => {
        setLoaded(false)
        try {
            const response = await axios.post('https://localhost/room/make', { uuid: context.userData.uuid});
            if (response.data) {
                if (response.data.state > 0) {
                    let userData = context.userData
                    userData.inRoom = true;
                    userData.roomAuthor = response.data.author.uuid == context.userData.uuid
                    userData.roomUuid = response.data.uuid
                    context.setUserData(userData)
                    navigate('/room');
                } else {
                    navigate('/home', { state: { message: 'Failed to create room' } });
                }
            } else {
                // well, we are on happy path
            }
        } catch (error) {
            navigate('/home', { state: { message: 'Failed to create room' } });
        }
    };

    const handleClickJoinRoomHandle = () => {
        setLoaded(false)
    };

    return (
        !loaded ? <Fallback/> :
        <>
            <Container className="p-5 mb-4 bg-light rounded-3">
                <h1 className="Header">React/Golang WebChat demo</h1>
            </Container>
            <h3>Hello, {context.userData.username} <Logout/></h3>
            <Form>
                <div className="card card-2">
                    <Container>
                        <Row>
                            <Col xs={3} className="text-center">
                                <Button type="button" variant="success" onClick={handleClickCreateRoom} className="btn-lg">Create&nbsp;room</Button>
                            </Col>
                            <Col xs={1} className="text-center">
                                <p className="some-pad-top">or</p>
                            </Col>
                            <Col xs={7} className="text-center">
                                <Form.Control placeholder="Enter room UUID..." className="some-margin-top" type="text" name="roomId" id="roomId"></Form.Control>
                            </Col>
                            <Col xs={1}>
                                <Button className="some-margin-top" type="button" variant="primary" onClick={handleClickJoinRoomHandle} >Join</Button>
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

export default HomePage