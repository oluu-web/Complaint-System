import React, { Fragment, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

const Complaint = () => {

 const { id } = useParams();
 const token = localStorage.getItem("token")

 const [complaint, setComplaint] = useState({
  id: id,
  matricNo: "",
  details: "",
 });
 const [isLoaded, setIsLoaded] = useState(false);
 const [error, setError] = useState(null);

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
   });
   setIsLoaded(true);
  })
  .catch((error) => {
   setIsLoaded(true);
   setError(error);
  })
 }, [token, id]);

 if(error) {
  return <div className='text-red-500'>Error: {error.message} </div>
 } else if (!isLoaded) {
  return <p>Loading...</p>
 } else {
  return (
   <Fragment>
    <p className='text-lg font-bold'>Matric Number: {complaint.matricNo}</p>
    <br />
    <br />

    <div className='bg-gray-100 p-4 rounded-lg'>
    <h2 className='text-xl font-bold'>Details</h2>
    <p>{complaint.details}</p>
    </div>
   </Fragment>
  )
 }
}

export default Complaint;