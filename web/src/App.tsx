import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Login from "./components/pages/login/login";
import ChatApp from "./components/pages/main/chatApp";
import ProtectedRoute from "./services/protectedRoute.service";



function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={<ProtectedRoute>{<ChatApp />}</ProtectedRoute>}
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
