import {ReactNode} from 'react';

interface PropsTypes {
    firstName: string;
    lastName: string;
    children: ReactNode;
}

const Profile = (props: PropsTypes) => {
    const {firstName, lastName} = props;
    return <h1>
        Name: {firstName} {lastName}
        {props.children}
    </h1>
};

export default Profile;