import { useContext, useEffect, useState } from 'react';
import { Button, Card, Col, Form, Row, ListGroup, ListGroupItem } from "react-bootstrap";
import { UserContext } from "../../App";
import { ChatClient } from "../../chat.client";
import { AuthRequest, ChatMessage, HistoryResponse} from "../../chat";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";

const ChatForm = () => {
    const { userData } = useContext(UserContext);
    const [messages, setMessages] = useState<ChatMessage[]>([]);

    let transport = new GrpcWebFetchTransport({
        baseUrl: `http://localhost:8080`
    });

    let client = new ChatClient(transport);

    const authRequest : AuthRequest = {
        uUID: userData.uuid,
        chatUUID: userData.roomUuid
    };

    const updateChat = async () => {
        const {response} = await client.getHistory(authRequest);
        response.messages.map((msg) => {
            setMessages(prevMessages => [...prevMessages, msg])
        })
    }

    useEffect(() => {
        updateChat().then().catch()
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