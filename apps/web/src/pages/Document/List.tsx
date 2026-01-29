import React, { useState, useEffect } from 'react';
import { Table, Button, Input, Space, Popconfirm, message, Tag } from 'antd';
import { SearchOutlined, DeleteOutlined, DownloadOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { getDocumentList, deleteDocument, downloadDocument, previewDocument, Document } from '../../services/document';
import { useNavigate } from 'react-router-dom';

const DocumentList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const navigate = useNavigate();

  const fetchDocuments = async () => {
    setLoading(true);
    try {
      const response = await getDocumentList({ page, page_size: pageSize, keyword });
      setDocuments(response.documents);
      setTotal(response.total);
    } catch (error) {
      message.error('获取文档列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDocuments();
  }, [page, pageSize]);

  const handleDelete = async (id: string) => {
    try {
      await deleteDocument(id);
      message.success('删除成功');
      fetchDocuments();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleDownload = (id: string, fileName: string) => {
    const url = downloadDocument(id);
    const link = document.createElement('a');
    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const handlePreview = (id: string) => {
    navigate(`/documents/${id}/preview`);
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  const getFileTypeTag = (fileType: string) => {
    const type = fileType.split('/')[0];
    const colorMap: Record<string, string> = {
      image: 'green',
      video: 'blue',
      audio: 'purple',
      application: 'orange',
      text: 'cyan',
    };
    return <Tag color={colorMap[type] || 'default'}>{type}</Tag>;
  };

  const columns: ColumnsType<Document> = [
    {
      title: '文件名',
      dataIndex: 'file_name',
      key: 'file_name',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'file_type',
      key: 'file_type',
      width: 120,
      render: (_, record) => getFileTypeTag(record.file_type),
    },
    {
      title: '大小',
      dataIndex: 'file_size',
      key: 'file_size',
      width: 120,
      render: (size) => formatFileSize(size),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (time) => new Date(time).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handlePreview(record.id)}
          >
            预览
          </Button>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record.id, record.file_name)}
          >
            下载
          </Button>
          <Popconfirm
            title="确定要删除这个文档吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Space>
          <Input
            placeholder="搜索文档"
            prefix={<SearchOutlined />}
            style={{ width: 300 }}
            onPressEnter={(e) => {
              setKeyword((e.target as HTMLInputElement).value);
              setPage(1);
            }}
            allowClear
          />
          <Button type="primary" onClick={() => fetchDocuments()}>
            搜索
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={documents}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />
    </div>
  );
};

export default DocumentList;
