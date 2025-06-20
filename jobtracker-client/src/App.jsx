import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Briefcase, Calendar, Building, FileText } from 'lucide-react';

const JobTracker = () => {
  const [jobs, setJobs] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingJob, setEditingJob] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [formData, setFormData] = useState({
    company: '',
    title: '',
    status: 'applied',
    applied_date: '', // Fixed typo: was 'appled_date'
    notes: ''
  });

  // API base URL - adjust this to your Go server
  const API_BASE = 'http://localhost:8080/jobs';

  // Fetch all jobs
  const fetchJobs = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch(API_BASE, {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        }
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      setJobs(data || []);
    } catch (error) {
      console.error('Error fetching jobs:', error);
      setError('Failed to fetch jobs: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  // Create or update job
  const saveJob = async () => {
    if (!formData.company.trim() || !formData.title.trim()) {
      alert('Please fill in company and title fields');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      if (editingJob) {
        // Update existing job
        const response = await fetch(`${API_BASE}/${editingJob.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const updatedJob = await response.json();
        setJobs(jobs.map(job => job.id === editingJob.id ? updatedJob : job));
      } else {
        // Create new job
        const response = await fetch(API_BASE, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const newJob = await response.json();
        setJobs([...jobs, newJob]);
      }

      closeModal();
    } catch (error) {
      console.error('Error saving job:', error);
      setError('Failed to save job: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  // Delete job
  const deleteJob = async (id) => {
    if (window.confirm('Are you sure you want to delete this job?')) {
      try {
        setLoading(true);
        setError(null);
        
        const response = await fetch(`${API_BASE}/${id}`, { method: 'DELETE' });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        setJobs(jobs.filter(job => job.id !== id));
      } catch (error) {
        console.error('Error deleting job:', error);
        setError('Failed to delete job: ' + error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  // Open modal for editing
  const openEditModal = (job) => {
    setEditingJob(job);
    setFormData({
      company: job.company,
      title: job.title,
      status: job.status,
      applied_date: job.applied_date || '',
      notes: job.notes
    });
    setIsModalOpen(true);
  };

  // Open modal for creating
  const openCreateModal = () => {
    setEditingJob(null);
    setFormData({
      company: '',
      title: '',
      status: 'applied',
      applied_date: new Date().toISOString().split('T')[0], // today's date
      notes: ''
    });
    setIsModalOpen(true);
  };

  // Close modal
  const closeModal = () => {
    setIsModalOpen(false);
    setEditingJob(null);
    setFormData({
      company: '',
      title: '',
      status: 'applied',
      applied_date: new Date().toISOString().split('T')[0], // today's date
      notes: ''
    });
    setError(null);
  };

  // Status colors
  const getStatusColor = (status) => {
    switch (status) {
      case 'applied': return 'bg-blue-100 text-blue-800';
      case 'interview': return 'bg-yellow-100 text-yellow-800';
      case 'offer': return 'bg-green-100 text-green-800';
      case 'rejected': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  // Load jobs on component mount
  useEffect(() => {
    fetchJobs();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-3">
            <Briefcase className="h-8 w-8 text-blue-600" />
            <h1 className="text-3xl font-bold text-gray-900">Job Tracker</h1>
          </div>
          <button
            onClick={openCreateModal}
            disabled={loading}
            className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            <Plus className="h-4 w-4" />
            Add Job
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
            <button 
              onClick={() => setError(null)}
              className="float-right text-red-700 hover:text-red-900"
            >
              ×
            </button>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="text-center py-4">
            <div className="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            <span className="ml-2">Loading...</span>
          </div>
        )}

        {/* Jobs Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {jobs.map((job) => (
            <div key={job.id} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-gray-900 mb-1">{job.title}</h3>
                  <div className="flex items-center gap-2 text-gray-600 mb-2">
                    <Building className="h-4 w-4" />
                    <span>{job.company}</span>
                  </div>
                  <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(job.status)}`}>
                    {job.status}
                  </span>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => openEditModal(job)}
                    disabled={loading}
                    className="p-2 text-gray-400 hover:text-blue-600 transition-colors disabled:opacity-50"
                  >
                    <Edit className="h-4 w-4" />
                  </button>
                  <button
                    onClick={() => deleteJob(job.id)}
                    disabled={loading}
                    className="p-2 text-gray-400 hover:text-red-600 transition-colors disabled:opacity-50"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>

              {(job.applied_date) && (
                <div className="flex items-center gap-2 text-sm text-gray-500 mb-3">
                  <Calendar className="h-4 w-4" />
                  <span>Applied: {job.applied_date}</span>
                </div>
              )}

              {job.notes && (
                <div className="flex items-start gap-2 text-sm text-gray-600">
                  <FileText className="h-4 w-4 mt-0.5" />
                  <p className="line-clamp-3">{job.notes}</p>
                </div>
              )}
            </div>
          ))}
        </div>

        {jobs.length === 0 && !loading && (
          <div className="text-center py-12">
            <Briefcase className="h-16 w-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No jobs yet</h3>
            <p className="text-gray-600 mb-4">Start tracking your job applications</p>
            <button
              onClick={openCreateModal}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
            >
              Add Your First Job
            </button>
          </div>
        )}
      </div>

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h2 className="text-xl font-semibold mb-4">
              {editingJob ? 'Edit Job' : 'Add New Job'}
            </h2>

            {error && (
              <div className="bg-red-100 border border-red-400 text-red-700 px-3 py-2 rounded mb-4 text-sm">
                {error}
              </div>
            )}

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Company *
                </label>
                <input
                  type="text"
                  value={formData.company}
                  onChange={(e) => setFormData({ ...formData, company: e.target.value })}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Job Title *
                </label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Status
                </label>
                <select
                  value={formData.status}
                  onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="applied">Applied</option>
                  <option value="interview">Interview</option>
                  <option value="offer">Offer</option>
                  <option value="rejected">Rejected</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Applied Date
                </label>
                <input
                  type="date"
                  value={formData.applied_date}
                  onChange={(e) => setFormData({ ...formData, applied_date: e.target.value })}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Notes
                </label>
                <textarea
                  value={formData.notes}
                  onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                  rows="3"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  onClick={saveJob}
                  disabled={loading}
                  className="flex-1 bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Saving...' : (editingJob ? 'Update' : 'Create')}
                </button>
                <button
                  onClick={closeModal}
                  disabled={loading}
                  className="flex-1 bg-gray-300 text-gray-700 py-2 rounded-md hover:bg-gray-400 transition-colors disabled:opacity-50"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default JobTracker;