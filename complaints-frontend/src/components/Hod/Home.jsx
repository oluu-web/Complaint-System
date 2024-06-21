import React, { useState, useEffect } from 'react';
import { Link } from "react-router-dom";
import '../Home.css'

const HODHome = () => {
 const [complaints, setComplaints] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1)
  const [complaintsPerPage, setComplaintsPerPage] = useState(10);
  const token = localStorage.getItem("token")
  // const userID = localStorage.getItem("userID")

  useEffect(() => {
   fetch("http://localhost:4000/hod-complaints", {
    headers: {
     Authorization: token,
    },
   })
   .then((response) => {
    if (response.status !== 200) {
     let err = new Error();
     err.message = "Invalid resposne code: " + response.status;
     throw err;
    }
    return response.json();
   })
   .then((json) => {
    setComplaints(json.complaints);
    setIsLoaded(true);
   })
   .catch((error) => {
    setIsLoaded(true);
    setError(error);
   })
  }, [token])

  const totalPages = complaints ? Math.ceil(complaints.length / complaintsPerPage) : 0;

  const handlePageChange = (pageNumber) => {
    setCurrentPage(pageNumber)
  }

  const handleComplaintsPerPageChange = (evt) => {
    setComplaintsPerPage(Number(evt.target.value));
    setCurrentPage(1);
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
   return <div> Error: {error.message}</div>
  } else if (!isLoaded) {
   return <p>Loading...</p>
  } else {
   const indexOfLastComplaint = currentPage * 10;
    const indexOfFirstComplaint = indexOfLastComplaint - 10;
    const currentComplaints = complaints ? complaints.slice(indexOfFirstComplaint, indexOfLastComplaint) : [];
  return (
   <div className="container mx-auto px-4 py-8">
   {complaints !== null ? (
    <table className = "table-auto w-full">
     <thead>
      <tr>
       <th className='border px-4 py-2'>Course concerned</th>
       <th className='border px-4 py-2'>Student Involved</th>
      </tr>
     </thead>
     <tbody>
      {currentComplaints.map((complaint) => (
       <tr className='border px-4 py-2' key = {complaint._id}>
       <Link to = {`complaint/${complaint._id}`}>
        <td className="border px-4 py-2">
         {complaint.course_concerned}
        </td>
        </Link>
        <td className="border px-4 py-2">
         {complaint.requesting_student}
        </td>
        {/* <td className="border px-4 py-2">
        <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">Approve</button>
      </td>
      <td className="border px-4 py-2">
        <button className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Decline</button>
      </td> */}
       </tr>
      ))}
     </tbody>
    </table>
   ) : (
    <div className='flex justify-center'>
     <h2>You have no re-evaluation requests at the moment</h2>
    </div>
   )}

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
  )
}
}
export default HODHome;