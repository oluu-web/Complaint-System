import React, { Fragment, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';

const Complaint = () => {
  const { id } = useParams();
  const token = localStorage.getItem("token");
  const navigate = useNavigate();

  const [complaint, setComplaint] = useState({
    id: id,
    matricNo: "",
    details: "",
    file_path: "",
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
          file_path: `http://localhost:4000/${json.complaint.file_path}`, // Update the file path
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

    axios.put(`http://localhost:4000/approved-by-lecturer/${id}`, {}, {
      headers: {
        Authorization: token,
      }
    })
      .then((res) => {
        console.log(res);
        navigate('/lecturer-dashboard');
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
      });
  };

  const handleDecline = (e) => {
    e.preventDefault();

    axios.put(`http://localhost:4000/decline/${id}`, {}, {
      headers: {
        Authorization: token,
      }
    })
      .then((res) => {
        console.log(res);
        navigate('/lecturer-dashboard');
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
            <h2 className="text-xl font-semibold mb-2">Details</h2>
            <p className="mb-4">{complaint.details}</p>
            {complaint.file_path && <img src={complaint.file_path} alt='complaint' className="max-w-full h-auto rounded-lg" />}
          </div>
          <div className="flex space-x-4">
            <button className="bg-green-500 text-white py-2 px-4 rounded hover:bg-green-600" onClick={handleAccept}>Accept</button>
            <button className="bg-red-500 text-white py-2 px-4 rounded hover:bg-red-600" onClick={handleDecline}>Decline</button>
          </div>
          {errorMessage && <p className="text-red-500 mt-4">{errorMessage}</p>}
        </div>
      </Fragment>
    );
  }
};

export default Complaint;
