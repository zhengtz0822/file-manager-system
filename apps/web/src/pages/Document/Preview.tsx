import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, Spin, Alert, message } from 'antd';
import { ArrowLeftOutlined, DownloadOutlined } from '@ant-design/icons';
import { getDocument, previewDocument, downloadDocument, Document } from '../../services/document';

const DocumentPreview: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [document, setDocument] = useState<Document | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchDocument();
  }, [id]);

  const fetchDocument = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const doc = await getDocument(id);
      setDocument(doc);
    } catch (err) {
      setError('获取文档信息失败');
      message.error('文档不存在或已被删除');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = () => {
    if (!document) return;
    const url = downloadDocument(document.id);
    const link = document.createElement('a');
    link.href = url;
    link.download = document.file_name;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const isImage = (fileType: string) => fileType.startsWith('image/');
  const isPDF = (fileType: string) => fileType === 'application/pdf';

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (error || !document) {
    return (
      <Alert
        message="错误"
        description={error || '文档不存在'}
        type="error"
        showIcon
        action={
          <Button size="small" onClick={() => navigate('/documents')}>
            返回列表
          </Button>
        }
      />
    );
  }

  return (
    <div>
      <Card
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/documents')}>
                返回
              </Button>
              <span>{document.file_name}</span>
            </Space>
            <Button
              type="primary"
              icon={<DownloadOutlined />}
              onClick={handleDownload}
            >
              下载
            </Button>
          </div>
        }
      >
        {isImage(document.file_type) || isPDF(document.file_type) ? (
          <div style={{ textAlign: 'center', background: '#f0f0f0', padding: '20px' }}>
            {isImage(document.file_type) ? (
              <img
                src={previewDocument(document.id)}
                alt={document.file_name}
                style={{ maxWidth: '100%', maxHeight: '800px' }}
              />
            ) : (
              <iframe
                src={previewDocument(document.id)}
                style={{ width: '100%', height: '800px', border: 'none' }}
              />
            )}
          </div>
        ) : (
          <Alert
            message="不支持预览"
            description={`文件类型 ${document.file_type} 暂不支持在线预览，请下载后查看`}
            type="warning"
            showIcon
          />
        )}

        <div style={{ marginTop: 16 }}>
          <Card type="inner" title="文件信息">
            <p><strong>文件名：</strong>{document.file_name}</p>
            <p><strong>文件大小：</strong>{(document.file_size / 1024 / 1024).toFixed(2)} MB</p>
            <p><strong>文件类型：</strong>{document.file_type}</p>
            <p><strong>上传时间：</strong>{new Date(document.created_at).toLocaleString('zh-CN')}</p>
          </Card>
        </div>
      </Card>
    </div>
  );
};

import { Space } from 'antd';

export default DocumentPreview;
