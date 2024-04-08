import { Route, BrowserRouter, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import StudentHome from "./components/Student/Home";
import ComplaintForm from "./components/Student/ComplaintForm";
import Test from "./components/test";
import Login from "./login";
import LecturerHome from "./components/Lecturer/Home";
import Complaint from "./components/Lecturer/ComplaintDetails";

function App() {
  return (
    <BrowserRouter>
      <Navbar />
        <Routes>
          <Route path="/" element={<Test />} />
          <Route path="/login" element={<Login />} />
          <Route path="/student-dashboard" element={<StudentHome />} />
          <Route path="/new-complaint" element={<ComplaintForm />} />
          <Route path="/lecturer-dashboard" element={<LecturerHome />} />
          <Route path="/lecturer-dashboard/complaint/:id" element={<Complaint />} />
        </Routes>
    </BrowserRouter>
  )
}

export default App;
