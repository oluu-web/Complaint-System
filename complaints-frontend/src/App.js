import { Route, BrowserRouter, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import StudentHome from "./components/Student/Home";
import ComplaintForm from "./components/Student/ComplaintForm";
import Test from "./components/test";
import Login from "./login";

function App() {
  return (
    <BrowserRouter>
      <Navbar />
        <Routes>
          <Route path="/" element={<Test />} />
          <Route path="/login" element={<Login />} />
          <Route path="/student-dashboard" element={<StudentHome />} />
          <Route path="/new-complaint" element={<ComplaintForm />} />
        </Routes>
    </BrowserRouter>
  )
}

export default App;
