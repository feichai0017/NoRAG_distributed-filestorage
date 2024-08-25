
import './index.css'
import React, { useState } from 'react';
import {queryAPI} from "@/api/files.jsx";

const QueryFile = () => {
    const [filehash, setFilehash] = useState('');
    const [results, setResults] = useState([]);
    const [error, setError] = useState('');

    const handleInputChange = (e) => {
        setFilehash(e.target.value);
    };

    const handleFormSubmit = async (e) => {
        e.preventDefault();

        try {
            const data = await queryAPI({ filehash: filehash });
            if (!data || data.error) {
                setError(data ? data.error : 'No data received');
                setResults([]);
            } else {
                setError('');
                setResults(data);
            }
        } catch (error) {
            setError('Error fetching data');
            setResults([]);
        }
    }

    return (
        <div className="searching-file">
            <div className="file-search">
            <h1>File Search</h1>
            <form onSubmit={handleFormSubmit}>
                <div className="input-box">
                    <input
                        type="text"
                        name="filehash"
                        placeholder="Enter filehash"
                        value={filehash}
                        onChange={handleInputChange}
                        required
                    />
                </div>
                <br />
                <div className="input-box button">
                    <input type="submit" value="Search" />
                </div>
                {error && <div className="error-message">{error}</div>}
                <div className="results">
                    <h2>Search Results</h2>
                    <ul id="results-list">
                        {results.map((result, index) => (
                            <li key={index}>
                                File Name: {result.file_name}, File Size: {result.file_size}, Location: {result.location}
                            </li>
                        ))}
                    </ul>
                </div>
            </form>
        </div>
        </div>
    );
};

export default QueryFile;