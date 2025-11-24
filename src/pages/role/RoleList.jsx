import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, message, Popconfirm } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getRoleList, createRole, updateRole, deleteRole } from '../../api/role'

const RoleList = () => {
    const [loading, setLoading] = useState(false)
    const [dataSource, setDataSource] = useState([])
    const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })
    const [modalVisible, setModalVisible] = useState(false)
    const [editingRole, setEditingRole] = useState(null)
    const [form] = Form.useForm()

    // 加载角色列表
    const loadData = async (page = 1, pageSize = 10) => {
        setLoading(true)
        try {
            const data = await getRoleList({ page, pageSize })
            setDataSource(data.list || [])
            setPagination({
                current: page,
                pageSize,
                total: data.total || 0,
            })
        } catch (error) {
            console.error('加载角色列表失败:', error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        loadData()
    }, [])

    // 表格列配置
    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '角色名称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '角色标识',
            dataIndex: 'code',
            key: 'code',
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
        },
        {
            title: '创建时间',
            dataIndex: 'created_at',
            key: 'created_at',
        },
        {
            title: '操作',
            key: 'action',
            width: 200,
            render: (_, record) => (
                <Space>
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定要删除这个角色吗？"
                        onConfirm={() => handleDelete(record.id)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    // 处理分页变化
    const handleTableChange = (newPagination) => {
        loadData(newPagination.current, newPagination.pageSize)
    }

    // 打开新增对话框
    const handleAdd = () => {
        setEditingRole(null)
        form.resetFields()
        setModalVisible(true)
    }

    // 打开编辑对话框
    const handleEdit = (record) => {
        setEditingRole(record)
        form.setFieldsValue(record)
        setModalVisible(true)
    }

    // 处理删除
    const handleDelete = async (id) => {
        try {
            await deleteRole(id)
            message.success('删除成功')
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            console.error('删除失败:', error)
        }
    }

    // 处理表单提交
    const handleSubmit = async () => {
        try {
            const values = await form.validateFields()

            if (editingRole) {
                await updateRole(editingRole.id, values)
                message.success('更新成功')
            } else {
                await createRole(values)
                message.success('创建成功')
            }

            setModalVisible(false)
            loadData(pagination.current, pagination.pageSize)
        } catch (error) {
            console.error('提交失败:', error)
        }
    }

    return (
        <div>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <h1>角色管理</h1>
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                    新增角色
                </Button>
            </div>

            <Table
                columns={columns}
                dataSource={dataSource}
                loading={loading}
                rowKey="id"
                pagination={pagination}
                onChange={handleTableChange}
            />

            <Modal
                title={editingRole ? '编辑角色' : '新增角色'}
                open={modalVisible}
                onOk={handleSubmit}
                onCancel={() => setModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="name"
                        label="角色名称"
                        rules={[{ required: true, message: '请输入角色名称' }]}
                    >
                        <Input placeholder="请输入角色名称" />
                    </Form.Item>

                    <Form.Item
                        name="code"
                        label="角色标识"
                        rules={[{ required: true, message: '请输入角色标识' }]}
                    >
                        <Input placeholder="请输入角色标识，如：admin" />
                    </Form.Item>

                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={4} placeholder="请输入角色描述" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default RoleList
