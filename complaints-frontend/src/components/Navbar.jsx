import React from 'react';
import { Link } from 'react-router-dom';

const Navbar = () => {
  return (
    <nav className="bg-gray-800 p-4">
      <div className="container mx-auto">
        <div className="flex items-center justify-between">
          <div>
            <Link to={"/"} className="text-white font-bold text-xl">Complaints System</Link>
          </div>
          <div>
            <Link to={"/complaints"} className="text-white mr-4">Complaints</Link>
            <Link to={"/new-complaint"} className="text-white">New Complaint</Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;