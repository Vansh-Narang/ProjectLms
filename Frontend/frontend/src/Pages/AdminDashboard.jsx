import React, { useState, useEffect } from "react";
import axios from "axios";
import "./Admin.css"
const AdminDashboard = () => {
    const token = localStorage.getItem("token");
    const [books, setBooks] = useState([]);
    const [requests, setRequests] = useState([]);
    // const [newBook, setNewBook] = useState({ title: "", author: "", year: "" });

    useEffect(() => {
        fetchBooks();
        fetchRequests();
    }, []);

    const [bookDetails, setBookDetails] = useState({
        isbn: "",
        title: "",
        authors: "",
        publisher: "",
        version: ""
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setBookDetails({
            ...bookDetails,
            [name]: name === "version" ? Number(value) : value, // Convert version to number
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch("http://localhost:8000/api/admin/add-book", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify(bookDetails)
            });
            const data = await response.json();
            alert("Book added successfully!");
        } catch (error) {
            console.error("Error adding book:", error);
            alert("Failed to add book");
        }
    };
    const fetchBooks = async () => {
        try {
            const response = await axios.get("/api/books");
            setBooks(response.data);
        } catch (error) {
            console.error("Error fetching books:", error);
        }
    };

    const fetchRequests = async () => {
        try {
            const response = await axios.get("http://localhost:8000/api/admin/list-requests",
                {
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${token}`
                    },
                }
            );
            console.log([response.data]);
            setRequests([response.data]);
        } catch (error) {
            console.error("Error fetching requests:", error);
        }
    };
    const removeBook = async (id) => {
        try {
            await axios.delete(`/api/books/${id}`);
            fetchBooks();
        } catch (error) {
            console.error("Error removing book:", error);
        }
    };

    const approveRequest = async (id) => {
        try {
            await axios.post(`/api/requests/${id}/approve`);
            fetchRequests();
        } catch (error) {
            console.error("Error approving request:", error);
        }
    };

    const rejectRequest = async (id) => {
        try {
            await axios.post(`/api/requests/${id}/reject`);
            fetchRequests();
        } catch (error) {
            console.error("Error rejecting request:", error);
        }
    };

    return (
        <div>
            <h1>Admin Dashboard</h1>

            <h2>Admin Dashboard - Add Book</h2>
            <form onSubmit={handleSubmit}>
                <input type="text" name="isbn" placeholder="ISBN" onChange={handleChange} required />
                <input type="text" name="title" placeholder="Title" onChange={handleChange} required />
                <input type="text" name="authors" placeholder="Authors" onChange={handleChange} required />
                <input type="text" name="publisher" placeholder="Publisher" onChange={handleChange} required />
                <input type="number" name="version" placeholder="Version" onChange={handleChange} required />
                <button type="submit">Add Book</button>
            </form>
            {/* <ul>
                {books.map((book) => (
                    <li key={book.id}>
                        {book.title} by {book.author} ({book.year})
                        <button onClick={() => removeBook(book.id)}>Remove</button>
                    </li>
                ))}
            </ul> */}

            <h2>Manage Requests</h2>
            <ul>
                {requests.map((req, index) => (
                    <li key={index}>
                        {req.user} requested {req.bookTitle}
                        <button onClick={() => approveRequest(req.id)}>Approve</button>
                        <button onClick={() => rejectRequest(req.id)}>Reject</button>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default AdminDashboard;
