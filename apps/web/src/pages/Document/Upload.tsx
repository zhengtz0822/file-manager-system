import React, { useState } from 'react';
import {
  Card,
  Upload,
  Button,
  Progress,
  message,
  Space,
  Alert,
  Descriptions,
} from 'antd';
import { InboxOutlined, CloudUploadOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import { initUpload, uploadChunk, completeUpload } from '../../services/document';

const { Dragger } = Upload;

const DocumentUpload: React.FC = () => {
  const [fileList, setFileList] = useState<any[]>([]);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [currentFile, setCurrentFile] = useState<any>(null);
  const [uploadedChunks, setUploadedChunks] = useState(0);
  const [totalChunks, setTotalChunks] = useState(0);

  const CHUNK_SIZE = 5 * 1024 * 1024; // 5MB

  const uploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    fileList: fileList,
    beforeUpload: (file) => {
      setFileList([file]);
      return false;
    },
    onRemove: () => {
      setFileList([]);
      setCurrentFile(null);
      setUploadProgress(0);
    },
  };

  const handleUpload = async () => {
    if (fileList.length === 0) {
      message.warning('请选择文件');
      return;
    }

    const file = fileList[0];
    setUploading(true);
    setCurrentFile(file);

    try {
      // 1. 初始化上传
      const initResp = await initUpload({
        file_name: file.name,
        file_size: file.size,
        chunk_size: CHUNK_SIZE,
      });

      setTotalChunks(initResp.total_chunks);

      // 2. 分片上传
      const chunks = Math.ceil(file.size / CHUNK_SIZE);
      let uploaded = 0;

      for (let i = 0; i < chunks; i++) {
        const start = i * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);

        await uploadChunk(initResp.upload_id, i + 1, chunk as any);

        uploaded++;
        setUploadedChunks(uploaded);
        setUploadProgress(Math.round((uploaded / chunks) * 100));
      }

      // 3. 完成上传
      const completeResp = await completeUpload({ upload_id: initResp.upload_id });

      message.success('上传成功！');
      setFileList([]);
      setCurrentFile(null);
      setUploadProgress(0);
      setUploadedChunks(0);

      // 可以在这里跳转到文档详情页
      // navigate(`/documents/${completeResp.document_id}`);
    } catch (error) {
      message.error('上传失败');
      console.error(error);
    } finally {
      setUploading(false);
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div>
      <Card title="上传文档" style={{ marginBottom: 16 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          <Dragger {...uploadProps} disabled={uploading}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p className="ant-upload-hint">
              支持单个文件上传，最大支持 5GB
            </p>
          </Dragger>

          {currentFile && (
            <Alert
              message="当前文件"
              description={
                <Descriptions size="small" column={1}>
                  <Descriptions.Item label="文件名">{currentFile.name}</Descriptions.Item>
                  <Descriptions.Item label="文件大小">
                    {formatFileSize(currentFile.size)}
                  </Descriptions.Item>
                  <Descriptions.Item label="分片数量">
                    {totalChunks > 0 ? `${totalChunks} 个` : '-'}
                  </Descriptions.Item>
                </Descriptions>
              }
              type="info"
            />
          )}

          {uploading && (
            <div>
              <div style={{ marginBottom: 8 }}>
                上传进度：{uploadedChunks}/{totalChunks} ({uploadProgress}%)
              </div>
              <Progress percent={uploadProgress} status="active" />
            </div>
          )}

          <Button
            type="primary"
            icon={<CloudUploadOutlined />}
            onClick={handleUpload}
            disabled={fileList.length === 0 || uploading}
            size="large"
            block
          >
            {uploading ? '上传中...' : '开始上传'}
          </Button>
        </Space>
      </Card>

      <Card title="上传说明">
        <Space direction="vertical">
          <div>• 支持的文件类型：图片、PDF、Office 文档、文本文件等</div>
          <div>• 单个文件最大 5GB，超过 10MB 会自动分片上传</div>
          <div>• 上传过程中请勿关闭页面</div>
        </Space>
      </Card>
    </div>
  );
};

export default DocumentUpload;
