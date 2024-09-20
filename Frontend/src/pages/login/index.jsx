// eslint-disable-next-line no-unused-vars
import React, { useState } from 'react'
import { motion } from 'framer-motion'
import { LockIcon, EyeIcon, EyeOffIcon, UserIcon } from 'lucide-react'
import { loginAPI } from "@/api/user.jsx"
import { setToken } from "@/utils/index.jsx"

export default function EnhancedLogin() {
    const [error, setError] = useState(null)
    const [showPassword, setShowPassword] = useState(false)
    const [isLoading, setIsLoading] = useState(false)

    const handleSubmit = async (event) => {
        event.preventDefault()
        setError(null)
        setIsLoading(true)
        const formData = new FormData(event.currentTarget)
        try {
            const response = await loginAPI(formData)
            console.log(response)
            if (response && response.code === 10000) {
                setToken(response.data.Token)
                localStorage.setItem('username', formData.get('username'))
                window.location.href = '/'
            } else {
                throw new Error('Unexpected response from server')
            }
        } catch (error) {
            console.error('Login error:', error)
            if (error.response) {
                switch (error.response.status) {
                    case 401:
                        setError('Invalid username or password')
                        break
                    case 500:
                        setError('An error occurred during login. Please try again later.')
                        break
                    default:
                        setError('An unexpected error occurred. Please try again.')
                }
            } else if (error.request) {
                setError('Unable to connect to the server. Please check your internet connection.')
            } else {
                setError('An unexpected error occurred. Please try again.')
            }
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-400 to-purple-500 p-8 flex items-center justify-center">
            <motion.div
                className="w-full max-w-md bg-white rounded-2xl shadow-2xl p-12 overflow-hidden relative"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
            >
                <motion.div
                    className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-blue-500 to-purple-600"
                    initial={{ scaleX: 0 }}
                    animate={{ scaleX: 1 }}
                    transition={{ duration: 0.5, delay: 0.2 }}
                />
                <div className="flex flex-col items-center mb-12">
                    <motion.div
                        className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-full p-4 mb-6"
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                    >
                        <LockIcon className="h-10 w-10 text-white" />
                    </motion.div>
                    <motion.h1
                        className="text-4xl font-bold text-gray-800 mb-2"
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.3 }}
                    >
                        Login
                    </motion.h1>
                    <motion.p
                        className="text-xl text-gray-600"
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.4 }}
                    >
                        Welcome back!
                    </motion.p>
                </div>
                {error && (
                    <motion.div
                        className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md mb-6"
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                    >
                        <p>{error}</p>
                    </motion.div>
                )}
                <form onSubmit={handleSubmit} className="space-y-6">
                    <div>
                        <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1 ml-2">
                            Username
                        </label>
                        <div className="relative">
                            <input
                                id="username"
                                name="username"
                                type="text"
                                autoComplete="username"
                                required
                                className="w-full px-4 py-2 pl-10 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                            <UserIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                        </div>
                    </div>
                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1 ml-2">
                            Password
                        </label>
                        <div className="relative">
                            <input
                                id="password"
                                name="password"
                                type={showPassword ? "text" : "password"}
                                autoComplete="current-password"
                                required
                                className="w-full px-4 py-2 pr-10 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                            <button
                                type="button"
                                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                                onClick={() => setShowPassword(!showPassword)}
                            >
                                {showPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                            </button>
                        </div>
                    </div>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center">
                            <input
                                id="remember-me"
                                name="remember-me"
                                type="checkbox"
                                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                            />
                            <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-600">
                                Remember me
                            </label>
                        </div>
                        <a href="#" className="text-sm font-medium text-blue-600 hover:text-purple-600 transition-colors">
                            Forgot your password?
                        </a>
                    </div>
                    <button
                        type="submit"
                        className={`w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-lg text-lg font-semibold transition-all duration-300 shadow-md hover:shadow-lg ${
                            isLoading ? 'opacity-50 cursor-not-allowed' : 'hover:from-blue-600 hover:to-purple-700'
                        }`}
                        disabled={isLoading}
                    >
                        {isLoading ? 'Logging in...' : 'Login'}
                    </button>
                </form>
                <motion.div
                    className="mt-8 text-center"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.5 }}
                >
                    <a href="/signup" className="text-blue-600 hover:text-purple-600 transition-colors text-lg">
                        {/* eslint-disable-next-line react/no-unescaped-entities */}
                        Don't have an account? Sign Up
                    </a>
                </motion.div>
            </motion.div>
        </div>
    )
}