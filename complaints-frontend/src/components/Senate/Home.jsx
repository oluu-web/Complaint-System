import React, { useState, useEffect, Fragment } from "react";
import { useNavigate } from "react-router-dom";
import '../Home.css';

const LecturerHome = () => {
  const [complaints, setComplaints] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [complaintsPerPage, setComplaintsPerPage] = useState(10);
  const navigate = useNavigate()
  const token = sessionStorage.getItem("token");

  useEffect(() => {
    fetch(`http://localhost:4000/senate-complaints`, {
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
        return response.json();
      })
      .then((json) => {
        if (json.complaints) {
          // Filter out complaints with the status "Approved By HOD"
          const filteredComplaints = json.complaints.filter(
            (complaint) => complaint.status === "Approved By HOD"
          );
          setComplaints(filteredComplaints);
        } else {
          setComplaints([]);
        }
        setIsLoaded(true);
      })
      .catch((error) => {
        setIsLoaded(true);
        setError(error);
      });
  }, [token]);

  const totalPages = complaints ? Math.ceil(complaints.length / complaintsPerPage) : 0;

  const handlePageChange = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  const handleComplaintsPerPageChange = (evt) => {
    setComplaintsPerPage(Number(evt.target.value));
    setCurrentPage(1);
  };

  const renderPaginationButtons = () => {
    const pageNumbers = [];
    for (let i = 1; i <= totalPages; i++) {
      pageNumbers.push(i);
    }

    return (
      <div className="flex justify-center mt-4">
        {pageNumbers.map((number) => (
          <button
            key={number}
            className={`mx-1 px-3 py-1 rounded ${currentPage === number ? "bg-blue-500 text-white" : "bg-gray-200"}`}
            onClick={() => handlePageChange(number)}
          >
            {number}
          </button>
        ))}
      </div>
    );
  };

  if (error) {
    return <div className="text-red-500 text-center mt-4">Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <p className="text-center mt-4">Loading...</p>;
  } else {
    const indexOfLastComplaint = currentPage * complaintsPerPage;
    const indexOfFirstComplaint = indexOfLastComplaint - complaintsPerPage;
    const currentComplaints = complaints.slice(indexOfFirstComplaint, indexOfLastComplaint);

    return (
      <Fragment>
        <div className="container mx-auto px-4 py-8">
          {complaints.length > 0 ? (
            <div>
              <table className="table-auto w-full">
                <thead>
                  <tr>
                    <th className="border px-4 py-2">Course Concerned</th>
                    <th className="border px-4 py-2">Student Involved</th>
                    <th className="border px-4 py-2">Assigned Lecturer</th>
                  </tr>
                </thead>
                <tbody>
                  {currentComplaints.map((complaint) => (
                    <tr key={complaint._id}
                    onClick={(e) => {
                        e.preventDefault();
                        navigate(`complaint/${complaint._id}`);
                      }}
                      style={{ cursor: "pointer" }}>
                      <td className="border px-4 py-2">{complaint.course_concerned}</td>
                      <td className="border px-4 py-2">{complaint.requesting_student}</td>
                      <td className="border px-4 py-2">{complaint.responding_lecturer}</td>
                    </tr>
                  ))}
                </tbody>
              </table>

              <div className="mt-8 flex justify-between items-center">
                {renderPaginationButtons()}
                <select
                  className="ml-4 px-2 py-1 rounded border"
                  value={complaintsPerPage}
                  onChange={handleComplaintsPerPageChange}
                >
                  <option value={10}>10</option>
                  <option value={50}>50</option>
                  <option value={100}>100</option>
                </select>
              </div>
            </div>
          ) : (
            <div className="flex justify-center">
              <h2>You have no re-evaluation requests at the moment</h2>
            </div>
          )}
        </div>
      </Fragment>
    );
  }
};

export default LecturerHome;
