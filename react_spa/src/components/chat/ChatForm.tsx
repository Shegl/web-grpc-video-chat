import {ChangeEvent, useContext, useEffect, useRef, useState} from 'react';
import { Button, Card, Col, Form, Row, ListGroup, ListGroupItem } from "react-bootstrap";
import { UserContext } from "../../App";
import { ChatClient } from "../../chat.client";
import { AuthRequest, ChatMessage } from "../../chat";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";

interface sendFormData {
    message: string;
}

const ChatForm = () => {
    const { userData } = useContext(UserContext);
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [chatInput, setChatInput] = useState<sendFormData>({message: ""});
    const lastChatElem = useRef<ListGroupItem>(null);

    let transport = new GrpcWebFetchTransport({
        baseUrl: `https://localhost`
    });

    let client = new ChatClient(transport);

    const authRequest : AuthRequest = {
        uUID: userData.uuid,
        chatUUID: userData.roomUuid
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setChatInput({
            ...chatInput,
            [e.target.name]: e.target.value,
        });
    };

    const handleSubmitChat = async (e) => {
        e.preventDefault();
        let chatMessageText: string = chatInput.message;
        if (chatMessageText.length > 0) {
            setChatInput({message: ""});
            const {response} = await client.sendMessage({
                msg: chatMessageText,
                authData: authRequest
            })
        }
    };

    const updateChatView = () => {
        setTimeout(function () {
           lastChatElem.current?.scrollIntoView({ behavior: "smooth" });
        }, 150);
    }

    const addMessage = (msg) => {
        setMessages((prevMessages) => {
            if (!prevMessages.some(messageItem => msg.uUID == messageItem.uUID)) {
                return [...prevMessages, msg]
            }
            return prevMessages
        });
    }

    const updateChat = async () => {
        const {response} = await client.getHistory(authRequest);
        response.messages.map((msg) => {
            addMessage(msg);
        });
        updateChatView();
    };

    const listenForMessages = async () => {
        let stream = client.listen(authRequest);
        for await (let message of stream.responses) {
            addMessage(message);
            updateChatView();
        }
    }

    useEffect(() => {
        listenForMessages();
        updateChat().then().catch();
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
                    <ListGroupItem key="lastelem" ref={lastChatElem} className="lastChatElem"></ListGroupItem>
                </ListGroup>
                <Card.Footer>
                    <Form onSubmit={handleSubmitChat}>
                        <Row>
                            <Col>
                                <Form.Control value={chatInput.message} type="text" name="message" id="message" onChange={handleChange}/>
                            </Col>
                            <Col xs={3}>
                                <Button type="submit" variant="primary inline">Send</Button>
                            </Col>
                        </Row>
                    </Form>
                </Card.Footer>
            </Card>
        </>
    );
};

export default ChatForm;