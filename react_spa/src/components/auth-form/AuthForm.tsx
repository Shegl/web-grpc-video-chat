import { useState, ChangeEvent, FormEvent, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Form, Col, Row, Button, Container} from "react-bootstrap";
import { useCookies } from "react-cookie";
import { UserContext } from "../../App";

interface AuthFormData {
    username: string;
}

const AuthForm = () => {
    const { setAuthenticated, userData, setUserData } = useContext(UserContext);
    const [_, setCookie] = useCookies(['userUuid']);
    const [formData, setFormData] = useState<AuthFormData>({
        username: '',
    });
    const navigate = useNavigate();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        try {
            const response = await axios.post('https://localhost/auth', formData);
            if (response.data.username && response.data.uuid) {

                userData.username = response.data.username;
                userData.uuid = response.data.uuid;

                setUserData(userData)
                setAuthenticated(true)
                setCookie('userUuid', response.data.uuid, { path: '/' });
                navigate('/home');
            } else {
                // well, we are on happy path
            }
        } catch (error) {
            console.error(error);
            console.error("Somethings is wrong, and backend seems to send errors");
        }
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    return (
        <Form onSubmit={handleSubmit}>
            <div className="card card-2">
                <Container>
                    <Row>
                        <Col xs={8}>
                            <Form.Label htmlFor="username">Nickname</Form.Label>
                        </Col>
                        <Col>
                        </Col>
                    </Row>
                    <Row>
                        <Col xs={8}>
                            <Form.Control type="text" name="username" id="username" onChange={handleChange}></Form.Control>
                        </Col>
                        <Col>
                            <Button type="submit" variant="primary">Set&nbsp;and&nbsp;proceed</Button>
                        </Col>
                    </Row>
                </Container>
            </div>
        </Form>
    );
};

export default AuthForm;