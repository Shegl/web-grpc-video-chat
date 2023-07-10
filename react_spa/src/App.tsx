import { useState } from 'react'
import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import {
    createBrowserRouter,
    RouterProvider,
    useLoaderData,
} from "react-router-dom"
import {Container} from "react-bootstrap";
import AuthForm from "./components/auth-form/AuthForm";

let router = createBrowserRouter([
    {
        path: "/",
        Component: Index
    },
    {
        path: "/rooms",
        Component: RoomsPage
    },
]);

function App() {
    return <RouterProvider router={router} fallbackElement={<Fallback />} />;
}

function Index() {
    return (
        <>
            <Container className="p-5 mb-4 bg-light rounded-3">
                <h1 className="Header">React/Golang WebChat demo</h1>
            </Container>
            <AuthForm></AuthForm>
            <p className="read-the-docs">
                Please provide your preferred name for the application and then press the button.
            </p>
        </>
    )
}

function RoomsPage() {
    const [userName, setUsername] = useState('')

    return (
        <>
            <h3>Hello, {userName}</h3>
            <div className="card">
                <button>
                    Create room
                </button>
            </div>
            <div className="card">
                <input type="text" name="roomId" />
                <button>
                    Join room
                </button>
            </div>
            <p className="read-the-docs">
                Create room or join existing one
            </p>
        </>
    )
}

export function Fallback() {
    return <p>Performing initial data load</p>;
}

export default App
