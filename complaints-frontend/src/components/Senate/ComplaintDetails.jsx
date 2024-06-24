import React, { Fragment, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

const Complaint = () => {
  const { id } = useParams();
  const token = sessionStorage.getItem("token");
  const [isAccepting, setIsAccepting] = useState(false);
  const [isDeclining, setIsDeclining] = useState(false);
  const navigate = useNavigate();

  const [complaint, setComplaint] = useState({
    id: id,
    matricNo: null,
    details: null,
    student_proof: null,
    lecturer_proof: null,
    reason: null,
  });
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    fetch(`http://localhost:4000/complaint/${id}`, {
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
        setComplaint({
          matricNo: json.complaint.requesting_student,
          details: json.complaint.request_details,
          student_proof: `http://localhost:4000/${json.complaint.student_proof}`,
          lecturer_proof: `http://localhost:4000/${json.complaint.lecturer_proof}`,
          reason: json.complaint.reason,
        });
        setIsLoaded(true);
      })
      .catch((error) => {
        setIsLoaded(true);
        setError(error);
      });
  }, [token, id]);

  const handleAccept = (e) => {
    e.preventDefault();
    setIsAccepting(true);

    axios.put(`http://localhost:4000/approved-by-senate/${id}`, {}, {
      headers: {
        Authorization: token,
      }
    })
      .then((res) => {
        console.log(res);
        navigate('/senate-dashboard');
      })
      .catch((err) => {
        if (err.response) {
          setErrorMessage("Failed to update complaint. Server responded with: " + JSON.stringify(err.response.data));
        } else if (err.request) {
          setErrorMessage("Failed to update complaint. No response received from the server.");
        } else {
          console.log('Error', err.message);
          setErrorMessage("Failed to submit complaint. Error: " + err.message);
        }
      })
      .finally(() => {
        setIsAccepting(false)
      });
  };

  const handleDecline = (e) => {
    e.preventDefault();
    setIsDeclining(true);

    axios.put(`http://localhost:4000/decline/${id}`, {}, {
      headers: {
        Authorization: token,
      }
    })
      .then((res) => {
        console.log(res);
        navigate('/senate-dashboard');
      })
      .catch((err) => {
        if (err.response) {
          setErrorMessage("Failed to update complaint. Server responded with: " + JSON.stringify(err.response.data));
        } else if (err.request) {
          setErrorMessage("Failed to update complaint. No response received from the server.");
        } else {
          console.log('Error', err.message);
          setErrorMessage("Failed to submit complaint. Error: " + err.message);
        }
      })
      .finally(() => {
        setIsDeclining(false);
      });
  };

  if (error) {
    return <div className='text-red-500'>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <p>Loading...</p>;
  } else {
    return (
      <Fragment>
        <div className="max-w-3xl mx-auto p-6 bg-white shadow-md rounded-lg">
          <h1 className="text-2xl font-bold mb-4">Complaint Details</h1>
          <p className="text-lg font-semibold mb-2">Matric Number: <span className="font-normal">{complaint.matricNo}</span></p>
          <div className="bg-gray-100 p-4 rounded-lg mb-4">
            <br />
            <h2 className="text-xl font-semibold mb-2"> Details</h2>
            <p className="mb-4">{complaint.details}</p>
            {complaint.student_proof && <img src={complaint.student_proof} alt='complaint' className="max-w-full h-auto rounded-lg" />}
            <h2 className="text-xl font-semibold mb-2">Approval Details</h2>
            <p><span className="font-semibold">Status</span>: {complaint.status}</p>
            <p className="mb-4"><span className="font-semibold">Reason</span>: {complaint.reason}</p>
            {complaint.status !== "Pending" ? (<img src={complaint.lecturer_proof} className="max-w-full h-auto rounded-lg" alt="lecturer proof"/>): (<div></div>)}
          </div>
          <div className="flex space-x-4">
            <button className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-600" onClick={handleAccept}>{isAccepting ? (
              <svg
                className="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                ></path>
              </svg>
            ) : (
              "Accept"
            )}</button>
            <button className="bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600" onClick={handleDecline}>
              {isDeclining ? (
              <svg
                className="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                ></path>
              </svg>
            ) : (
              "Decline"
            )}
            </button>
          </div>
          {errorMessage && <p className="text-red-500 mt-4">{errorMessage}</p>}
        </div>
      </Fragment>
    );
  }
};

export default Complaint;
