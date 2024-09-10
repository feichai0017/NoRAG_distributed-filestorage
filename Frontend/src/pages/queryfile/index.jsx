// eslint-disable-next-line no-unused-vars
import React, { useState } from 'react';
import { Search, FileText, AlertCircle } from 'lucide-react';
import { queryAPI } from "@/api/files.jsx";

const EnhancedQueryFile = () => {
    const [filehash, setFilehash] = useState('');
    const [results, setResults] = useState([]);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleInputChange = (e) => {
        setFilehash(e.target.value);
    };

    const handleFormSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');
        setResults([]);

        try {
            const data = await queryAPI({ filehash: filehash });
            if (!data || data.error) {
                setError(data ? data.error : 'No data received');
            } else {
                setResults(data);
            }
        } catch (error) {
            setError('Error fetching data');
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center p-4">
            <div className="w-full max-w-2xl bg-white rounded-2xl shadow-xl overflow-hidden transition-all duration-300 ease-in-out transform hover:scale-105">
                <div className="p-8">
                    <h1 className="text-3xl font-bold mb-6 text-center bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-purple-500">File Search</h1>
                    <form onSubmit={handleFormSubmit} className="space-y-6">
                        <div className="relative">
                            <input
                                type="text"
                                name="filehash"
                                placeholder="Enter filehash"
                                value={filehash}
                                onChange={handleInputChange}
                                required
                                className="w-full px-4 py-3 rounded-lg border-2 border-gray-200 focus:border-blue-500 focus:ring focus:ring-blue-200 focus:ring-opacity-50 transition-all duration-300 ease-in-out pl-12"
                            />
                            <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400" />
                        </div>
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="w-full px-6 py-3 border border-transparent text-base font-medium rounded-lg text-white bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-300 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {isLoading ? 'Searching...' : 'Search'}
                        </button>
                    </form>

                    {error && (
                        <div className="mt-6 p-4 bg-red-100 border-l-4 border-red-500 rounded-r-lg">
                            <div className="flex items-center">
                                <AlertCircle className="h-6 w-6 text-red-500 mr-2" />
                                <p className="text-red-700">{error}</p>
                            </div>
                        </div>
                    )}

                    {results.length > 0 && (
                        <div className="mt-8">
                            <h2 className="text-2xl font-semibold mb-4 text-gray-800">Search Results</h2>
                            <ul className="space-y-4">
                                {results.map((result, index) => (
                                    <li key={index} className="bg-gray-50 rounded-lg p-4 shadow transition-all duration-300 ease-in-out hover:shadow-md">
                                        <div className="flex items-start">
                                            <FileText className="h-6 w-6 text-blue-500 mr-3 mt-1" />
                                            <div>
                                                <p className="font-medium text-gray-800">{result.file_name}</p>
                                                <p className="text-sm text-gray-600">Size: {result.file_size}</p>
                                                <p className="text-sm text-gray-600">Location: {result.location}</p>
                                            </div>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default EnhancedQueryFile;