import {createPortal} from "react-dom";
import { useRecoilState } from "recoil";
import styles from './modal.module.css';
import GlobalState from "../GlobalState";
import { useState } from "react";

const Modal = () => {
   const [globalUserName, setGlobalUserName] = useRecoilState(GlobalState.userNameState);
   const [userName, setUserName] = useState('');
   const handleChange = (event: React.FormEvent<HTMLInputElement>) => {
      setUserName(event.currentTarget.value as string);
  }
  const handleSubmit = (event: React.FormEvent) => {
      event.preventDefault();
      setGlobalUserName(userName);
      setUserName('');
  }
  const handleCancel = (event: React.FormEvent) => {
   event.preventDefault();
      setUserName('');
  }
     return createPortal(
      <div className={styles.modalForm}>
      <form onSubmit={handleSubmit}>
          <input value={userName} onChange={handleChange} className={styles.userNameInput} placeholder="Username"></input>
          <div className={styles.modalButtons}>
          <button type="submit" className={styles.loginButton}>Login</button>
          <button className={styles.cancelButton} onClick={handleCancel}>Cancel</button>
          </div>
         <div>GlobalUserName: {globalUserName}</div>
      </form>
  </div>, document.getElementById('modal') as Element
     );
}

export default Modal;