import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, theme } from 'antd';
import LoginPage from './pages/Login';
import MainLayout from './components/Layout/MainLayout';
import DocumentList from './pages/Document/List';
import DocumentUpload from './pages/Document/Upload';
import DocumentPreview from './pages/Document/Preview';
import { useAuth } from './hooks/useAuth';

const App: React.FC = () => {
  const { isAuthenticated } = useAuth();

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={isAuthenticated ? <Navigate to="/" replace /> : <LoginPage />}
        />
        <Route
          path="/"
          element={
            isAuthenticated ? (
              <MainLayout>
                <Routes>
                  <Route index element={<Navigate to="/documents" replace />} />
                  <Route path="documents" element={<DocumentList />} />
                  <Route path="upload" element={<DocumentUpload />} />
                  <Route path="documents/:id/preview" element={<DocumentPreview />} />
                </Routes>
              </MainLayout>
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
