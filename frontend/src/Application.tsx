import React from "react";
import ReactDom from "react-dom/client";
import {BrowserRouter, Routes, Route} from 'react-router-dom';
import About from "./pages/About";
import Home from "./pages/Home";

const domContainer = document.querySelector('#application');
const root = ReactDom.createRoot(domContainer!);
root.render(<BrowserRouter>
    <Routes>
      <Route index element={<Home />} />
      <Route path="/about" element={<About />} />
    </Routes>
  </BrowserRouter>);