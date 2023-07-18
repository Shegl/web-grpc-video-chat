import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import {
    BrowserRouter,
    Routes, Route
} from "react-router-dom"
import {Button, Container} from "react-bootstrap";
import AuthForm from "./components/auth-form/AuthForm";
import HomePage from "./components/rooms-actions/HomePage";
import { useCookies } from "react-cookie";
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import RoomPage from "./components/rooms-actions/RoomPage";

const defaultUserContext = {
    username: 'Anonymous',
    uuid: '',
    inRoom: false,
    roomAuthor: false,
    roomUuid: "",
};

export const UserContext = createContext({
    authenticated: false,
    userData: defaultUserContext,
    setAuthenticated: (value: boolean) => { console.log(value) },
    setUserData: (obj: any) => { console.log(obj) },
});


function App() {
    const [authenticated, setAuthenticated] = useState(false)
    const [userData, setUserData] = useState(defaultUserContext)

    const value = useMemo(
        () => ({ authenticated, setAuthenticated, userData, setUserData }),
        [authenticated, userData]
    );

    return (<BrowserRouter>
                <Routes>
                    <Route
                        path="/"
                        element={
                            <UserContext.Provider value={value}>
                                {useMemo(() => (
                                    <StartPage/>
                                ), [])}
                            </UserContext.Provider>
                        }
                    />
                    <Route
                        path="/home"
                        element={
                            <UserContext.Provider value={value}>
                                {useMemo(() => (
                                    <HomePage/>
                                ), [])}
                            </UserContext.Provider>
                        }
                    />
                    <Route
                        path="/room"
                        element={
                            <UserContext.Provider value={value}>
                                {useMemo(() => (
                                    <RoomPage/>
                                ), [])}
                            </UserContext.Provider>
                        }
                    />
                </Routes>
            </BrowserRouter>
    );
}

export const useAuth = (
    page: string,
    redirect: boolean,
    cookies: any,
    setCookie: any,
    navigate: any,
    context: any,
    setLoaded: any,
) => {
    const { authenticated, setAuthenticated, setUserData } = context;
    if (authenticated) {
        if (redirect) {
            navigate(page);
        } else {
            setLoaded(true);
        }
    } else {
        let userUUID = cookies.userUuid;
        if (userUUID) {
            axios.post('https://localhost/check', {uuid: userUUID}).then(
                (response) => {
                    if (response.status == 200) {
                        setAuthenticated(true);
                        setUserData({
                            username: response.data.username,
                            uuid: response.data.uuid,
                        })
                        if (redirect) {
                            navigate(page);
                        } else {
                            setLoaded(true);
                        }
                    } else {
                        setAuthenticated(false);
                        navigate("/");
                        setLoaded(true);
                    }
                }
            ).catch(
                (reason) => {
                    console.log(reason)
                    setCookie('userUuid', '', { path: '/' });
                    navigate('/')
                }
            )
        } else {
            if (!redirect && page != '/') {
                navigate('/');
            } else {
                setLoaded(true);
            }
        }
    }
}

function StartPage() {
    const [loaded, setLoaded] = useState(false);
    const [cookies, setCookie] = useCookies(['userUuid']);
    const navigate = useNavigate();
    const context = useContext(UserContext);

    useEffect(() => {
        useAuth("/home", true, cookies, setCookie, navigate, context, setLoaded);
    });

    return (
        !loaded ? <Fallback/> :
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

export function Fallback() {

    return <p>Loading...</p>;
}

export function Logout() {
    const [_, setCookie] = useCookies(['userUuid']);
    const navigate = useNavigate();
    const context = useContext(UserContext);

    const handleLogout = () => {
        setCookie('userUuid', '', { path: '/' });
        context.setUserData(defaultUserContext);
        context.setAuthenticated(false);

        axios.post('https://localhost/logout', {uuid: context.userData.uuid}).then().catch()

        navigate('/');
    }

    return  <Button type="button" variant="link" className="btn-sm" onClick={handleLogout}>[ Logout ]</Button>;
}


export default App
