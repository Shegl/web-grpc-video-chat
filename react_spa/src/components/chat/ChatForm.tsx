import { useContext, useEffect, useState } from 'react';
import { Button, Card, Col, Form, Row, ListGroup, ListGroupItem } from "react-bootstrap";
import { UserContext } from "../../App";
import { ChatClient } from "../../chat.client";
import { AuthRequest, ChatMessage } from "../../chat";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";

const ChatForm = () => {
    const { userData } = useContext(UserContext);
    const [messages, setMessages] = useState<ChatMessage[]>([]);

    let transport = new GrpcWebFetchTransport({
        baseUrl: `https://localhost`
    });

    let client = new ChatClient(transport);

    const authRequest : AuthRequest = {
        uUID: userData.uuid,
        chatUUID: userData.roomUuid
    };

    const updateChat = async () => {
        const {response} = await client.getHistory(authRequest);
        response.messages.map((msg) => {
            setMessages((prevMessages) => {
                if (!prevMessages.some(messageItem => msg.uUID == messageItem.uUID)) {
                    return [...prevMessages, msg]
                }
                return prevMessages
            })
        });
    }

    useEffect(() => {
        updateChat().then().catch()
    }, [userData]);

    return (
        <>
            <Card>
                <Card.Header>
                    Chat
                </Card.Header>
                <ListGroup className="messages-window" variant="flush">
                    {messages.map(message => (
                        <ListGroupItem key={message.uUID}>
                            <strong>{ message.userName }{ message.userUUID == userData.uuid ? "(You)" : "" }</strong>
                        : {message.msg}</ListGroupItem>
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