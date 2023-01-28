import {ReactNode} from 'react';

interface PropsTypes {
    firstName: string;
    lastName: string;
    children: ReactNode;
}

const Profile = (props: PropsTypes) => {
    return <h1>
        Name: {props.firstName} {props.lastName}
        {props.children}
    </h1>
};

export default Profile;