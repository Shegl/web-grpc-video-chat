import { useState, ChangeEvent, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Form, Col, Row, Button, Container} from "react-bootstrap";

interface AuthFormData {
    username: string;
}

const AuthForm = () => {
    const [formData, setFormData] = useState<AuthFormData>({
        username: '',
    });

    const navigate = useNavigate();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        try {
            const response = await axios.post('http://localhost:3000/auth', formData);
            if (response.data.success) {
                navigate('/rooms');
            } else {
                // well, we are on happy path
            }
        } catch (error) {
            navigate('/rooms');
            console.error("Somethings is wrong, and backend seems to send errors")
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
            <div className="card">
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