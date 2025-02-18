// src/components/Dashboard.tsx
import React from 'react';
import { useAuth } from '../context/AuthContext';

const Dashboard: React.FC = () => {
    const { user } = useAuth();

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <div className="space-y-6">
                <div>
                    <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
                    <p className="mt-1 text-gray-600">
                        Welcome to your dashboard, {user?.username}!
                    </p>
                </div>

                <div className="border-t border-gray-200 pt-6">
                    <h3 className="text-lg font-medium text-gray-900">Your Profile</h3>
                    <dl className="mt-4 space-y-4">
                        <div className="grid grid-cols-3 gap-4">
                            <dt className="text-sm font-medium text-gray-500">Username</dt>
                            <dd className="text-sm text-gray-900 col-span-2">{user?.username}</dd>
                        </div>
                        <div className="grid grid-cols-3 gap-4">
                            <dt className="text-sm font-medium text-gray-500">Email</dt>
                            <dd className="text-sm text-gray-900 col-span-2">{user?.email}</dd>
                        </div>
                        <div className="grid grid-cols-3 gap-4">
                            <dt className="text-sm font-medium text-gray-500">Account ID</dt>
                            <dd className="text-sm text-gray-900 col-span-2">{user?.id}</dd>
                        </div>
                    </dl>
                </div>

                <div className="border-t border-gray-200 pt-6">
                    <h3 className="text-lg font-medium text-gray-900">Quick Actions</h3>
                    <div className="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
                        <button className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            Update Profile
                        </button>
                        <button className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                            View Settings
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;