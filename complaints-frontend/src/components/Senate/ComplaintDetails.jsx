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
  const [errorMessage, setErrorMessage] = useState("")

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

  axios.put(`http://localhost:4000/approved-by-senate/${id}`, {}, {
    headers: {
      Authorization: token,
    }
  })
  .then((res) => {
    console.log(res)
    navigate('/senate-dashboard')
  })
  .catch((err) => {
    if (err.response) {
      setErrorMessage("Failed to update complaint. Server responded with: " + JSON.stringify(err.response.data))
    } else if (err.request)  {
      setErrorMessage("Failed toupdate complaint. No response recieved from the server.");
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
    navigate('/lecturer-dashboard')
  })
  .catch(err => {
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
    return <div className='text-red-500'>Error: {error.message} </div>;
  } else if (!isLoaded) {
    return <p>Loading...</p>;
  } else {
    return (
      <Fragment>
        <p className='text-lg font-bold m-30'>Matric Number: {complaint.matricNo}</p>
        <br />
        <br />

        <div className='bg-gray-100 p-4 rounded-lg'>
          <h2 className='text-xl font-bold'>Details</h2>
          <p>{complaint.details}</p>
          <br />
          {complaint.file_path !== '' && <img src={complaint.file_path} alt='complaint'/>}
        </div>
        <div className='mt-4'>
          <button className='bg-green-500 text-white p-2 rounded mr-2' onClick={handleAccept}>Accept</button>
          <button className='bg-red-500 text-white p-2 rounded' onClick={handleDecline}>Decline</button>
        </div>
        {errorMessage && <p className='text-red-500'>{errorMessage}</p>}
      </Fragment>
    );
  }
};

export default Complaint;
