import { Route, BrowserRouter, Routes } from "react-router-dom";
import StudentHome from "./components/Student/Home";
import ComplaintForm from "./components/Student/ComplaintForm";
import Login from "./login";
import LecturerHome from "./components/Lecturer/Home";
import Complaint from "./components/Lecturer/ComplaintDetails";
import Navbar from "./components/Navbar";
import HODHome from "./components/Hod/Home";

function App() {
  return (
    <BrowserRouter>
    <Navbar />
        <Routes>
          <Route path="/" element={<Login />} />
          {/* <Route path="/login" element={<Login />} /> */}
          <Route path="/student-dashboard" element={<StudentHome />} />
          <Route path="/new-complaint" element={<ComplaintForm />} />
          <Route path="/lecturer-dashboard" element={<LecturerHome />} />
          <Route path="/lecturer-dashboard/complaint/:id" element={<Complaint />} />
          <Route path="/hod-dashboard" element={<HODHome />} />
        </Routes>
    </BrowserRouter>
  )
}

export default App;
