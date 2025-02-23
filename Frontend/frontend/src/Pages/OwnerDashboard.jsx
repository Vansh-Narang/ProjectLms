import React, { useState } from "react";
import "./OwnerDashboard.css";

const OwnerDashboard = () => {
    const [showLibraryModal, setShowLibraryModal] = useState(false);
    const [showAdminModal, setShowAdminModal] = useState(false);
    const [libraryName, setLibraryName] = useState("");
    const [adminName, setAdminName] = useState("");
    const [adminEmail, setAdminEmail] = useState("");
    const [adminContact, setAdminContact] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const handleLibraryClick = () => {
        setShowLibraryModal(true);
    };

    const closeLibraryModal = () => {
        setShowLibraryModal(false);
        setLibraryName("");
        setError("");
    };

    const handleAdminClick = () => {
        setShowAdminModal(true);
    };

    const closeAdminModal = () => {
        setShowAdminModal(false);
        setAdminName("");
        setAdminEmail("");
        setAdminContact("");
        setError("");
    };

    // API Call: Create Library
    const handleCreateLibrary = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError("");

        try {
            const response = await fetch("http://localhost:8000/api/library/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({ name: libraryName }),
            });

            const data = await response.json();
            if (!response.ok) throw new Error(data.error || "Failed to create library");

            alert("Library created successfully!");
            closeLibraryModal();
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // API Call: Create Admin
    const handleCreateAdmin = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError("");

        try {
            const response = await fetch("http://localhost:8000/api/library/create-admin", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({
                    name: adminName,
                    email: adminEmail,
                    contactNumber: adminContact,
                }),
            });

            const data = await response.json();
            if (!response.ok) throw new Error(data.error || "Failed to create admin");

            alert("Admin created successfully!");
            closeAdminModal();
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="owner-dashboard">
            <h1 className="dashboard-title">Owner Dashboard</h1>
            <div className="button-container">
                <button className="create-library" onClick={handleLibraryClick}>Create Library</button>
                <button className="create-admin" onClick={handleAdminClick}>Create Admin</button>
            </div>

            {/* Library Modal */}
            {showLibraryModal && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Create Library</h2>
                        {error && <p className="error">{error}</p>}
                        <form onSubmit={handleCreateLibrary}>
                            <label>Library Name:</label>
                            <input
                                type="text"
                                placeholder="Enter library name"
                                value={libraryName}
                                onChange={(e) => setLibraryName(e.target.value)}
                                required
                            />
                            <button type="submit" disabled={loading}>{loading ? "Creating..." : "Submit"}</button>
                            <button type="button" onClick={closeLibraryModal}>Close</button>
                        </form>
                    </div>
                </div>
            )}

            {/* Admin Modal */}
            {showAdminModal && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Create Admin</h2>
                        {error && <p className="error">{error}</p>}
                        <form onSubmit={handleCreateAdmin}>
                            <label>Name:</label>
                            <input
                                type="text"
                                placeholder="Enter admin name"
                                value={adminName}
                                onChange={(e) => setAdminName(e.target.value)}
                                required
                            />
                            <label>Email:</label>
                            <input
                                type="email"
                                placeholder="Enter admin email"
                                value={adminEmail}
                                onChange={(e) => setAdminEmail(e.target.value)}
                                required
                            />
                            <label>Contact Number:</label>
                            <input
                                type="text"
                                placeholder="Enter contact number"
                                value={adminContact}
                                onChange={(e) => setAdminContact(e.target.value)}
                                required
                            />
                            <button type="submit" disabled={loading}>{loading ? "Creating..." : "Submit"}</button>
                            <button type="button" onClick={closeAdminModal}>Close</button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OwnerDashboard;
