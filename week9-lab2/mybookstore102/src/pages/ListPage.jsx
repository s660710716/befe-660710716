import React, { useState, useEffect } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';

const ListPage = () => {
  // กำหนด State สำหรับจัดการข้อมูล
  const [featuredBooks, setFeaturedBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
      useEffect(() => {
    // Check authentication
    const isAuthenticated = localStorage.getItem('isAdminAuthenticated');
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [navigate]);

  useEffect(() => {
    const fetchBooks = async () => {
      try {
        setLoading(true);
        
        // เรียก API เพื่อดึงข้อมูลหนังสือ
        const response = await fetch('http://localhost:8080/api/v1/books');

        if (!response.ok) {
          throw new Error('Failed to fetch books');
        }

        const data = await response.json();

        setFeaturedBooks(data);
        setError(null);
        
      } catch (err) {
        setError(err.message);
        console.error('Error fetching books:', err);
        
      } finally {
        setLoading(false);
      }
    };

    // เรียกใช้ฟังก์ชันดึงข้อมูล
    fetchBooks();
  }, []); // [] = dependency array ว่าง = รันครั้งเดียว

  // กรณีกำลังโหลดข้อมูล
  if (loading) {
    return (
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div className="text-center py-8 col-span-full">
          Loading...
        </div>
      </div>
    );
  }

  // กรณีเกิดข้อผิดพลาด
  if (error) {
    return (
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div className="text-center py-8 col-span-full text-red-600">
          Error: {error}
        </div>
      </div>
    );
  }

  // กรณีแสดงผลข้อมูลปกติ
  return (
    <div>
      <div className="bg-white text-black rounded-xl shadow-2xl overflow-hidden">
        <div>
            <buttom type="add" className="px-4 py-0 bg-green-400 text-viridian-600"><NavLink to = "/store-manager/add-book">เพิ่มหนังสือ</NavLink></buttom>
        </div>
      <table className="w-full border border-gray-200 rounded-2xl shadow-sm overflow-hidden">
        <thead className="bg-gradient-to-r from-green-300 to-green-200 text-gray-800">
          <tr>
            <th className="px-4 py-3 text-left font-semibold">ID</th>
            <th className="px-4 py-3 text-left font-semibold">ชื่อ</th>
            <th className="px-4 py-3 text-left font-semibold">ผู้แต่ง</th>
            <th className="px-4 py-3 text-left font-semibold">ปี</th>
            <th className="px-4 py-3 text-left font-semibold">ราคา</th>
            <th className="px-4 py-3 text-left font-semibold">ประเภท</th>
            <th className="px-4 py-3 text-left font-semibold">คะแนน</th>
            <th className="px-4 py-3"></th>
            <th className="px-4 py-3"></th>
          </tr>
        </thead>

        <tbody className="divide-y divide-gray-200 bg-white">
          {featuredBooks.map((book) => (
            <tr
              key={book.id}
              className="hover:bg-green-50 transition-all duration-150 ease-in-out"
            >
              <td className="px-4 py-3 text-gray-700">{book.id}</td>
              <td className="px-4 py-3 text-gray-800 font-medium">{book.title}</td>
              <td className="px-4 py-3 text-gray-600">{book.author}</td>
              <td className="px-4 py-3 text-gray-600">{book.year}</td>
              <td className="px-4 py-3 text-gray-600">{book.price}</td>
              <td className="px-4 py-3 text-gray-600">{book.category}</td>
              <td className="px-4 py-3 text-yellow-600 font-semibold">{book.rating}</td>
              <td className="px-4 py-3">
                <NavLink
                  to="/books/edit"
                  className="inline-block px-3 py-1 bg-blue-500 text-white text-sm rounded-xl hover:bg-blue-600 transition-colors"
                >
                  แก้ไข
                </NavLink>
              </td>
              <td className="px-4 py-3">
                <NavLink
                  to="/books/delete"
                  className="inline-block px-3 py-1 bg-red-500 text-white text-sm rounded-xl hover:bg-red-600 transition-colors"
                >
                  ลบ
                </NavLink>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

        
      </div>
    </div>
  );
};

export default ListPage;