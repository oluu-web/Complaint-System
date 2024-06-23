import { Route, BrowserRouter, Routes, useLocation } from "react-router-dom";
import StudentHome from "./components/Student/Home";
import ComplaintForm from "./components/Student/ComplaintForm";
import Login from "./login";
import Navbar from "./components/Navbar";
import LecturerHome from "./components/Lecturer/Home";
import LecturerComplaint from "./components/Lecturer/ComplaintDetails";
import HODHome from "./components/Hod/Home";
import HODComplaint from "./components/Hod/ComplaintDetails";
import SenateHome from "./components/Senate/Home";
import SenateComplaint from "./components/Senate/ComplaintDetails";

function App() {
  const location = useLocation();

  // List of paths where Navbar should not be shown
  const noNavbarPaths = ["/"];

  return (
    <>
      {!noNavbarPaths.includes(location.pathname) && <Navbar />}
      <Routes>
        <Route path="/" element={<Login />} />
        {/* <Route path="/login" element={<Login />} /> */}
        <Route path="/student-dashboard" element={<StudentHome />} />
        <Route path="/new-complaint" element={<ComplaintForm />} />
        <Route path="/lecturer-dashboard" element={<LecturerHome />} />
        <Route path="/lecturer-dashboard/complaint/:id" element={<LecturerComplaint />} />
        <Route path="/hod-dashboard" element={<HODHome />} />
        <Route path="/hod-dashboard/complaint/:id" element={<HODComplaint />} />
        <Route path="/senate-dashboard" element={<SenateHome />} />
        <Route path="/senate-dashboard/complaint/:id" element={<SenateComplaint />} />
      </Routes>
    </>
  );
}

function AppWrapper() {
  return (
    <BrowserRouter>
      <App />
    </BrowserRouter>
  );
}

export default AppWrapper;
