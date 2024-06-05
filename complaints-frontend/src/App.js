import { Route, BrowserRouter, Routes } from "react-router-dom";
import StudentHome from "./components/Student/Home";
import ComplaintForm from "./components/Student/ComplaintForm";
import Login from "./login";
import Navbar from "./components/Navbar";
import LecturerHome from "./components/Lecturer/Home";
import LecturerComplaint from "./components/Lecturer/ComplaintDetails";
import HODHome from "./components/Hod/Home";
import HODComplaint from "./components/Hod/ComplaintDetails"
import SenateHome from "./components/Senate/Home";
import SenateComplaint from "./components/Senate/ComplaintDetails"

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
          <Route path="/lecturer-dashboard/complaint/:id" element={<LecturerComplaint />} />
          <Route path="/hod-dashboard" element={<HODHome />} />
          <Route path="/hod-dashboard/complaint/:id" element = {<HODComplaint />} />
          <Route path="/senate-dashboard" element={<SenateHome />} />
          <Route path="/senate-dashboard/complaint/:id" element = {<SenateComplaint />} />
        </Routes>
    </BrowserRouter>
  )
}

export default App;
