import React, { useState } from 'react';

const QueryFile = () => {
    const [filehash, setFilehash] = useState('');
    const [results, setResults] = useState([]);

    const handleInputChange = (e) => {
        setFilehash(e.target.value);
    };

    const handleFormSubmit = (e) => {
        e.preventDefault();

        const params = new URLSearchParams({ filehash });

        fetch('/file/meta', {
            method: 'POST',
            body: params,
        })
            .then((response) => response.json())
            .then((data) => {
                setResults([{
                    file_name: data.file_name,
                    file_size: data.file_size,
                    location: data.location,
                }]);
            })
            .catch((error) => console.error('Error:', error));
    };

    return (
        <div className="container">
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
    );
};

export default QueryFile;