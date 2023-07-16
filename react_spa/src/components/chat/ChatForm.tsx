import { useContext, useEffect, useState } from 'react';
import { Button, Card, Col, Form, Row, ListGroup, ListGroupItem } from "react-bootstrap";
import { UserContext } from "../../App";
import { ChatClient } from "../../ChatServiceClientPb";
import {AuthRequest, ChatMessage} from "../../chat_pb";

const ChatForm = () => {
    const client = new ChatClient("http://localhost:8080", null, null);
    const { userData } = useContext(UserContext);
    const [messages, setMessages] = useState<ChatMessage[]>([]);

    const authRequest = new AuthRequest();
    console.log(authRequest.toString());
    authRequest.setUuid(userData.uuid);
    authRequest.setChatuuid(userData.roomUuid);

    useEffect(() => {
        client.getHistory(authRequest, null, (err, response) => {
            if (err) return console.log(err);
            response.getMessagesList().map(msg => {
                setMessages(prevMessages => [...prevMessages, msg]);
            })
        });
    });

    return (
        <>
            <Card>
                <Card.Header>
                    Chat
                </Card.Header>
                <ListGroup className="messages-window" variant="flush">
                    {messages.map(message => (
                        <ListGroupItem>{message.toString()}</ListGroupItem>
                    ))}
                </ListGroup>
                <Card.Footer>
                    <Row>
                        <Col><Form.Control type="text" name="chat-input" id="chat-input"></Form.Control></Col>
                        <Col xs={3}><Button type="button" variant="primary inline">Send</Button></Col>
                    </Row>
                </Card.Footer>
            </Card>
        </>
    );
};

export default ChatForm;