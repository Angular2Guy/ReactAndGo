
interface ListItemProps {
    name: string;
}
/*
const ListItem = (props: ListItemProps) => {
    return (<div>{props.name}</div>);
}
*/
const ListItem = ({name}: ListItemProps) => {
    return (<div>{name}</div>);
}

export default ListItem;