import {useContext, useEffect, useState} from 'react';
import {Fallback, Logout, useAuth, UserContext} from "../../App";
import { Button, Col, Container, Form, Row } from "react-bootstrap";
import { useCookies } from "react-cookie";
import { useNavigate } from 'react-router-dom';

function HomePage() {
    const [loaded, setLoaded] = useState(false);
    const context = useContext(UserContext);
    const [cookies, setCookie] = useCookies(['userUuid']);
    const navigate = useNavigate();

    useEffect(() => {
        useAuth("/home", false, cookies, setCookie, navigate, context, setLoaded);
    });

    return (
        !loaded ? <Fallback/> :
        <>
            <Container className="p-5 mb-4 bg-light rounded-3">
                <h1 className="Header">React/Golang WebChat demo</h1>
            </Container>
            <h3>Hello, {context.userData.username} <Logout/></h3>
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

export default HomePage