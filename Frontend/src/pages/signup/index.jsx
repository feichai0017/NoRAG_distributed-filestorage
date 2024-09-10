// eslint-disable-next-line no-unused-vars
import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { LockIcon, CheckCircleIcon, EyeIcon, EyeOffIcon, UserIcon, MailIcon } from 'lucide-react'
import { signUpAPI } from "@/api/user.jsx"

export default function EnhancedSignUp() {
    const [isSuccess, setIsSuccess] = useState(false)
    const [error, setError] = useState(null)
    const [showPassword, setShowPassword] = useState(false)
    const [formErrors, setFormErrors] = useState({})
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: ''
    })

    const validateForm = () => {
        const errors = {}
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

        if (!formData.username) errors.username = 'Username is required'
        if (!emailRegex.test(formData.email)) errors.email = 'Invalid email address'
        if (formData.password.length < 8) errors.password = 'Password must be at least 8 characters'

        return errors
    }

    const handleInputChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({ ...prev, [name]: value }))
    }

    const handleSubmit = async (event) => {
        event.preventDefault()
        setError(null)
        setFormErrors({})
        const errors = validateForm()

        if (Object.keys(errors).length === 0) {
            try {
                const response = await signUpAPI(formData)
                if (response && response.status === 201) {
                    setIsSuccess(true)
                } else {
                    throw new Error(response.data || 'An unexpected error occurred')
                }
            } catch (error) {
                console.error('Sign up failed:', error)
                if (error.response && error.response.status === 409) {
                    setError('Username already exists. Please choose a different username.')
                } else {
                    setError('Sign up failed. Please try again. ' + (error.message || ''))
                }
            }
        } else {
            setFormErrors(errors)
        }
    }

    const SuccessModal = () => {
        useEffect(() => {
            const timer = setTimeout(() => {
                window.location.href = '/login'
            }, 3000)
            return () => clearTimeout(timer)
        }, [])

        return (
            <motion.div
                className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
            >
                <motion.div
                    className="bg-white rounded-2xl p-8 max-w-md w-full shadow-2xl"
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    exit={{ scale: 0.9, opacity: 0 }}
                >
                    <motion.div
                        initial={{ scale: 0 }}
                        animate={{ scale: 1 }}
                        transition={{ delay: 0.2, type: "spring", stiffness: 200 }}
                    >
                        <CheckCircleIcon className="h-24 w-24 text-green-500 mx-auto mb-6" />
                    </motion.div>
                    <h2 className="text-4xl font-bold text-gray-800 mb-4 text-center">Welcome Aboard!</h2>
                    <p className="text-xl text-gray-600 mb-6 text-center">Your account has been successfully created.</p>
                    <motion.div
                        className="w-full bg-gradient-to-r from-blue-500 to-purple-600 h-2 rounded-full overflow-hidden"
                        initial={{ width: 0 }}
                        animate={{ width: "100%" }}
                        transition={{ duration: 3 }}
                    />
                    <p className="text-sm text-gray-500 mt-4 text-center">Redirecting to login page...</p>
                </motion.div>
            </motion.div>
        )
    }

    const inputClasses = "w-full px-4 py-3 text-lg border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-300 bg-gray-50"
    const labelClasses = "block text-sm font-medium text-gray-700 mb-1 ml-2"

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-400 to-purple-500 p-8 flex items-center justify-center">
            <AnimatePresence>
                {isSuccess && <SuccessModal />}
            </AnimatePresence>
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
                        Sign Up
                    </motion.h1>
                    <motion.p
                        className="text-xl text-gray-600"
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.4 }}
                    >
                        Join our community today
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
                        <label htmlFor="username" className={labelClasses}>Username</label>
                        <div className="relative">
                            <input
                                id="username"
                                name="username"
                                value={formData.username}
                                onChange={handleInputChange}
                                required
                                className={inputClasses}
                            />
                            <UserIcon className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                        </div>
                        {formErrors.username && <p className="text-red-500 text-sm mt-1 ml-2">{formErrors.username}</p>}
                    </div>
                    <div>
                        <label htmlFor="email" className={labelClasses}>Email Address</label>
                        <div className="relative">
                            <input
                                id="email"
                                name="email"
                                type="email"
                                value={formData.email}
                                onChange={handleInputChange}
                                required
                                className={inputClasses}
                            />
                            <MailIcon className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                        </div>
                        {formErrors.email && <p className="text-red-500 text-sm mt-1 ml-2">{formErrors.email}</p>}
                    </div>
                    <div>
                        <label htmlFor="password" className={labelClasses}>Password</label>
                        <div className="relative">
                            <input
                                id="password"
                                name="password"
                                type={showPassword ? "text" : "password"}
                                value={formData.password}
                                onChange={handleInputChange}
                                required
                                className={inputClasses}
                            />
                            <button
                                type="button"
                                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                                onClick={() => setShowPassword(!showPassword)}
                            >
                                {showPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                            </button>
                        </div>
                        {formErrors.password && <p className="text-red-500 text-sm mt-1 ml-2">{formErrors.password}</p>}
                    </div>
                    <motion.button
                        type="submit"
                        className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-lg text-lg font-semibold mt-8 transition-all duration-300 shadow-md hover:shadow-lg"
                        whileHover={{ scale: 1.02 }}
                        whileTap={{ scale: 0.98 }}
                    >
                        Sign Up
                    </motion.button>
                </form>
                <motion.div
                    className="mt-8 text-center"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.5 }}
                >
                    <a href="/login" className="text-blue-600 hover:text-purple-600 transition-colors text-lg">
                        Already have an account? Login
                    </a>
                </motion.div>
            </motion.div>
        </div>
    )
}