import React from 'react';

const Navbar = () => {
  return (
    <nav className="bg-purple-950 p-4">
      <div className="container mx-auto">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-white font-bold text-xl">CU Revalidation System</p>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;