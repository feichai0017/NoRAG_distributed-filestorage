
import './index.css'
import React, { useState } from 'react';

const Upload = () => {
    const [file, setFile] = useState(null);
    const [fileName, setFileName] = useState('');

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        setFile(selectedFile);
        setFileName(selectedFile.name);
    };

    const handleButtonClick = () => {
        document.getElementById('fileID').click();
    };

    return (
        <div className="upload-file">
            <div className="card">
                <h3>Upload Files</h3>
                <div className="drop_box">
                    <header>
                        <h4>Select File here</h4>
                    </header>
                    <p>Files Supported: PDF, TEXT, DOC, DOCX</p>
                    <input type="file" id="fileID" style={{ display: 'none' }} onChange={handleFileChange} />
                    <button className="btn" onClick={handleButtonClick}>Choose File</button>
                </div>
                {file && (
                    <form action="" method="post" encType="multipart/form-data">
                        <div className="form">
                            <h4>{fileName}</h4>
                            <input type="text" name="username" placeholder="Enter your name" />
                            <input type="hidden" name="filename" value={fileName} />
                            <input type="file" name="file" style={{ display: 'none' }} id="hiddenFile" />
                            <button type="submit" className="btn">Upload</button>
                        </div>
                    </form>
                )}
            </div>
        </div>
    );
};

export default Upload;