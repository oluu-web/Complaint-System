import React, { useState, useEffect} from "react";
import { useNavigate } from "react-router-dom";
import '../Home.css'
// import Navbar from "../Navbar";
const StudentHome = () => {
  const [complaints, setComplaints] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1)
  const [complaintsPerPage, setComplaintsPerPage] = useState(10);
  const navigate = useNavigate();
  const token = localStorage.getItem("token")
  const userID = localStorage.getItem("userID")
  useEffect(() => {
  fetch(`http://localhost:4000/complaints/${userID}`, {
    headers: {
      Authorization: token,
    },
  })
  .then((response) => {
    if (response.status !== 200) {
    let err = new Error();
          err.message = "Invalid response code: " + response.status;
          throw err;
  }
  return response.json()
  })
  .then((json) => {
    setComplaints(json.complaints);
    setIsLoaded(true);
  })
  .catch((error) => {
    setIsLoaded(true);
    setError(error);
  })
  }, [token, userID])
  const totalPages = complaints ? Math.ceil(complaints.length / complaintsPerPage) : 0;

  const handlePageChange = (pageNumber) => {
    setCurrentPage(pageNumber)
  }

  const handleComplaintsPerPageChange = (evt) => {
    setComplaintsPerPage(Number(evt.target.value));
    setCurrentPage(1);
  };
  const handleNewComplaint = () => {
    navigate('/new-complaint');
  };

  const renderPaginationButtons = () => {
    const pageNumbers = []
    for (let i = 1; i <= totalPages; i++) {
      pageNumbers.push(i);
    }

    return (
      <div className="pagination">
        {pageNumbers.map((number) => (
          <button
          key={number}
          className={currentPage === number ? "active" : ""}
          onClick={() => handlePageChange(number)}
        >{number}</button>
        ))}
      </div>
    )
  }

  if (error) {
    return <div> Error: {error.message} </div>
  } else if (!isLoaded) {
    return <p>Loading...</p>
  } else {
    const indexOfLastComplaint = currentPage * 10;
    const indexOfFirstComplaint = indexOfLastComplaint - 10;
    const currentComplaints = complaints ? complaints.slice(indexOfFirstComplaint, indexOfLastComplaint) : [];
  return (
    <div>
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">My Complaints</h1>
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          onClick={handleNewComplaint}
        >
          New Complaint
        </button>
      </div>
      {complaints !== null ? (
        <table className="table-auto w-full">
          <thead>
            <tr>
              <th className="border px-4 py-2">Course Concerned</th>
              <th className="border px-4 py-2">Assigned Lecturer</th>
              <th className="border px-4 py-2">Status</th>
            </tr>
          </thead>
          <tbody>
            {currentComplaints.map((complaint) => (
              <tr key={complaint._id}>
                <td className="border px-4 py-2">{complaint.course_concerned}</td>
                <td className="border px-4 py-2">{complaint.responding_lecturer}</td>
                <td className="border px-4 py-2">{complaint.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <div className="flex justify-center">
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            onClick={handleNewComplaint}
          >
            New Complaint
          </button>
        </div>
      )}
    </div>

    <div className='pagination-container'>
          {renderPaginationButtons()}
          <select
            // className="pagination-select"
            value={complaintsPerPage}
            onChange={handleComplaintsPerPageChange}
          >
            <option value={10}>10</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
        </div>
    </div>
  );
}
}

export default StudentHome