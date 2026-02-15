<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="名称">
          <el-input v-model.trim="params.groupName" style="width: 100px;" clearable placeholder="名称" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model.trim="params.remark" style="width: 100px;" clearable placeholder="描述" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="同步状态">
          <el-select v-model.trim="params.syncState" style="width: 110px;" clearable placeholder="同步状态" @change="search" @clear="search">
            <el-option label="已同步" value="1" />
            <el-option label="未同步" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <!-- <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="resetData">重置</el-button>
        </el-form-item> -->
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-upload2" type="success" @click="batchSync">批量同步</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="previewLoading" icon="el-icon-view" type="info" plain @click="syncPreview">同步预览</el-button>
        </el-form-item>
        <el-form-item>
          <el-button icon="el-icon-setting" type="primary" plain @click="$router.push('/personnel/user')">目录快速配置</el-button>
        </el-form-item>
        <el-form-item>
          <el-dropdown trigger="click" @command="handleGroupColumnCommand">
            <el-button type="default" plain size="mini" icon="el-icon-s-operation">列设置</el-button>
            <el-dropdown-menu slot="dropdown" class="column-setting-dropdown">
              <el-dropdown-item command="reset"><i class="el-icon-refresh-left" /> 重置为默认</el-dropdown-item>
              <el-dropdown-item divided disabled>显示列</el-dropdown-item>
              <el-dropdown-item v-for="col in defaultGroupColumns" :key="col.prop" :command="col.prop">
                <el-checkbox :value="columnConfig.visible[col.prop]" @click.native.prevent>{{ col.label }}</el-checkbox>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-form-item>
        <el-form-item>
          <el-tag size="small" type="info">目录类型：{{ directoryTypeText }}</el-tag>
        </el-form-item>
        <br>
        <el-form-item v-if="syncConfig.ldapEnableSync">
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncOpenLdapDepts">同步原ldap部门</el-button>
        </el-form-item>
        <el-form-item v-if="syncConfig.dingTalkEnableSync">
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncDingTalkDepts">同步钉钉部门</el-button>
        </el-form-item>
        <el-form-item v-if="syncConfig.feiShuEnableSync">
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncFeiShuDepts">同步飞书部门</el-button>
        </el-form-item>
        <el-form-item v-if="syncConfig.weComEnableSync">
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncWeComDepts">同步企业微信部门</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :default-expand-all="true" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" row-key="ID" :data="infoTableData" border stripe style="width: 100%" @selection-change="handleSelectionChange" @header-dragend="handleGroupTableHeaderDragend">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column
          v-for="col in visibleGroupColumns"
          :key="col.prop"
          show-overflow-tooltip
          sortable
          :prop="col.prop"
          :label="col.label"
          :width="col.width"
          :min-width="col.minWidth"
        />
        <el-table-column fixed="right" label="操作" align="center" width="220">
          <template #default="scope">
            <el-tooltip v-if="scope.row.groupType != 'ou' && scope.row.groupName != 'root'" content="添加" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-setting" circle type="info" @click="addUp(scope.row)" />
            </el-tooltip>
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip v-if="scope.row.syncState === 2" class="delete-popover" content="同步" effect="dark" placement="top">
              <el-popconfirm title="确定同步吗？" @onConfirm="singleSync(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-upload2" circle type="success" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <!-- 新增 -->
      <el-dialog :title="dialogFormTitle" :visible.sync="updateLoading">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="名称" prop="groupName">
            <el-input v-model.trim="dialogFormData.groupName" placeholder="名称(拼音)" />
          </el-form-item>
          <el-form-item label="分组类型" prop="groupType">
            <el-select v-model.trim="dialogFormData.groupType" placeholder="建议仅第一层为ou，如果不确定，就用cn" style="width:100%">
              <el-option label="cn[分组]" value="cn" />
              <el-option label="ou[组织]" value="ou" />
            </el-select>
          </el-form-item>
          <el-form-item label="上级分组" prop="parentId">
            <treeselect
              v-model="dialogFormData.parentId"
              :options="treeselectData"
              :normalizer="normalizer"
              placeholder="请选择上级分组"
              @input="treeselectInput"
            />
          </el-form-item>
          <el-form-item label="描述" prop="remark">
            <el-input v-model.trim="dialogFormData.remark" type="textarea" placeholder="描述" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
      <!-- 同步预览（Dry Run）结果 -->
      <el-dialog title="同步预览" :visible.sync="previewDialogVisible" width="520px" append-to-body @close="previewResult = null">
        <div v-if="previewResult" class="sync-preview-body">
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="将新增到 LDAP">{{ previewResult.addCount }} 个部门</el-descriptions-item>
            <el-descriptions-item label="LDAP 中已存在（将更新）">{{ previewResult.updateCount }} 个部门</el-descriptions-item>
          </el-descriptions>
          <template v-if="(previewResult.addList && previewResult.addList.length) || (previewResult.updateList && previewResult.updateList.length)">
            <p v-if="previewResult.addList && previewResult.addList.length" class="preview-list"><strong>新增：</strong>{{ previewResult.addList.join('、') }}</p>
            <p v-if="previewResult.updateList && previewResult.updateList.length" class="preview-list"><strong>已存在：</strong>{{ previewResult.updateList.join('、') }}</p>
          </template>
          <p class="preview-tip">以上为预览结果，未执行实际同步。点击「执行同步」将正式同步到 LDAP。</p>
        </div>
        <div slot="footer" class="dialog-footer">
          <el-button size="small" @click="previewDialogVisible = false">关 闭</el-button>
          <el-button size="small" type="success" :disabled="!previewResult || (previewResult.addCount === 0 && previewResult.updateCount === 0)" @click="confirmSyncFromPreview">执行同步</el-button>
        </div>
      </el-dialog>

      <!-- 编辑 -->
      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="名称" prop="groupName">
            <el-input v-model.trim="dialogFormData.groupName" :disabled="true" placeholder="名称" />
          </el-form-item>
          <el-form-item label="描述" prop="remark">
            <el-input v-model.trim="dialogFormData.remark" type="textarea" placeholder="描述" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { getGroupTree, groupAdd, groupUpdate, groupDel, syncDingTalkDeptsApi, syncWeComDeptsApi, syncFeiShuDeptsApi, syncOpenLdapDeptsApi, syncSqlGroups, syncSqlGroupsPreview } from '@/api/personnel/group'
import { getConfig } from '@/api/system/base'
import { loadTableColumnConfig, saveTableColumnConfig, defaultGroupColumns, STORAGE_KEY_GROUP_TABLE } from '@/utils/tableColumnSettings'
import { Message } from 'element-ui'

export default {
  name: 'Group',
  components: {
    Treeselect
  },
  filters: {
    methodTagFilter(val) {
      if (val === 'GET') {
        return ''
      } else if (val === 'POST') {
        return 'success'
      } else {
        return 'info'
      }
    }
  },
  data() {
    return {
      // 查询参数
      params: {
        groupName: undefined,
        remark: undefined,
        syncState: undefined,
        pageNum: 1,
        pageSize: 1000// 平常百姓人家应该不会有这么多数据吧,后台限制最大单次获取1000条
      },
      // 表格数据
      tableData: [],
      infoTableData: [],
      total: 0,
      loading: false,
      // 上级目录数据
      treeselectData: [],
      treeselectValue: 0,
      updateLoading: false, // 新增
      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        ID: '',
        groupName: '',
        parentId: 0,
        syncState: 1,
        groupType: '',
        remark: ''
      },
      dialogFormRules: {

        groupName: [
          { required: true, message: '请输入所属类别', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        groupType: [
          { required: true, message: '请输入分组类型', trigger: 'blur' },
          { min: 1, max: 50, message: 'ou、cn或者其他', trigger: 'blur' }
        ],
        parentId: [
          { required: true, message: '请选择父级', trigger: 'blur' },
          { validator: (rule, value, callBack) => {
            if (value >= 0) {
              callBack()
            } else {
              callBack('请选择有效的部门')
            }
          } }
        ],
        remark: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 100 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      dialogTransfer: '', // 穿梭框头部
      dialogTransferVisible: false,

      transParams: {
        groupId: '',
        nickname: ''
      },
      renderFunc(h, option) {
        return <span>{option.key} - {option.label}</span>
      },
      userArrInfo: [], // 初始人员列表数据
      data: [], // 转化后人员列表数据
      value3: [], // 右侧默认人员列表数据
      userId: [], // 送到后台 -> 勾选的数据code数组
      ui: {
        submitLoading: false
      },
      statusTrans: '',

      // 同步配置
      syncConfig: {
        ldapEnableSync: false,
        dingTalkEnableSync: false,
        feiShuEnableSync: false,
        weComEnableSync: false,
        directoryType: 'openldap'
      },

      // 表格列设置（显示/隐藏、列宽），从 localStorage 恢复
      columnConfig: (() => {
        const c = loadTableColumnConfig(STORAGE_KEY_GROUP_TABLE, defaultGroupColumns)
        return { visible: c.visible, widths: c.widths }
      })(),

      previewLoading: false,
      previewDialogVisible: false,
      previewResult: null
    }
  },
  created() {
    this.getTableData()
    this.getSyncConfig()
  },
  computed: {
    directoryTypeText() {
      const t = (this.syncConfig.directoryType || '').toLowerCase()
      if (t === 'ad') {
        return 'Windows AD'
      }
      return 'OpenLDAP'
    },
    visibleGroupColumns() {
      return defaultGroupColumns
        .filter(c => this.columnConfig.visible[c.prop] !== false)
        .map(c => ({
          ...c,
          width: this.columnConfig.widths[c.prop] || c.width || undefined,
          minWidth: c.minWidth
        }))
    }
  },
  methods: {
    toggleGroupColumnVisible(prop) {
      this.columnConfig.visible[prop] = !this.columnConfig.visible[prop]
      saveTableColumnConfig(STORAGE_KEY_GROUP_TABLE, this.columnConfig)
    },
    handleGroupColumnCommand(cmd) {
      if (cmd === 'reset') this.resetGroupColumnSettings()
      else this.toggleGroupColumnVisible(cmd)
    },
    resetGroupColumnSettings() {
      this.columnConfig.visible = {}
      this.columnConfig.widths = {}
      defaultGroupColumns.forEach(c => {
        this.$set(this.columnConfig.visible, c.prop, c.visible !== false)
        if (c.width) this.$set(this.columnConfig.widths, c.prop, c.width)
      })
      saveTableColumnConfig(STORAGE_KEY_GROUP_TABLE, this.columnConfig)
      Message.success('已重置为默认列')
    },
    handleGroupTableHeaderDragend(newWidth, oldWidth, column) {
      if (column && column.property) {
        this.$set(this.columnConfig.widths, column.property, newWidth)
        saveTableColumnConfig(STORAGE_KEY_GROUP_TABLE, this.columnConfig)
      }
    },
    // 获取同步配置
    async getSyncConfig() {
      try {
        const { data } = await getConfig()
        this.syncConfig = {
          ...this.syncConfig,
          ...data
        }
      } catch (error) {
        console.error('获取同步配置失败:', error)
      }
    },
    // // 查询
    search() {
      // 初始化表格数据
      this.infoTableData = JSON.parse(JSON.stringify(this.tableData))
      this.infoTableData = this.deal(this.infoTableData, node => node.groupName.includes(this.params.groupName) || node.remark.includes(this.params.remark) || node.syncState.toString().includes(this.params.syncState))
    },
    resetData() {
      this.infoTableData = JSON.parse(JSON.stringify(this.tableData))
    },
    // 页面数据过滤
    deal(nodes, predicate) {
      // 如果已经没有节点了，结束递归
      if (!(nodes && nodes.length)) {
        return []
      }
      const newChildren = []
      for (const node of nodes) {
        if (predicate(node)) {
          // 如果节点符合条件，直接加入新的节点集
          newChildren.push(node)
          node.children = this.deal(node.children, predicate)
        } else {
          // 如果当前节点不符合条件，递归过滤子节点，
          // 把符合条件的子节点提升上来，并入新节点集
          newChildren.push(...this.deal(node.children, predicate))
        }
      }
      return newChildren
    },
    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getGroupTree(this.params)
        this.tableData = data
        this.infoTableData = JSON.parse(JSON.stringify(data))
        this.treeselectData = [{ ID: 0, groupName: '顶级类目', children: data }]
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增分组'
      this.updateLoading = true // 新增的展示
      this.dialogType = 'create'
    },
    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.groupName = row.groupName
      this.dialogFormData.remark = row.remark
      this.dialogFormTitle = '修改分组'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },
    // 穿梭框
    addUp(row) {
      this.dialogTransfer = '用户管理'
      this.dialogTransferVisible = true
      this.transParams.groupId = row.ID
      this.transParams.nickname = row.remark
      this.$router.push({ path: '/userList', query: row })
    },

    // 判断结果
    judgeResult(res) {
      if (res.code === 0) {
        Message({
          showClose: true,
          message: '操作成功',
          type: 'success'
        })
      }
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              await groupAdd(this.dialogFormData).then(res => {
                this.judgeResult(res)
              })
            } else {
              await groupUpdate(this.dialogFormData).then(res => {
                this.judgeResult(res)
              })
            }
          } finally {
            this.submitLoading = false
          }
          this.resetForm()
          this.getTableData()
        } else {
          Message({
            showClose: true,
            message: '表单校验失败',
            type: 'warn'
          })
          return false
        }
      })
    },

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.updateLoading = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {

        groupName: '',
        remark: ''
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const groupIds = []
        this.multipleSelection.forEach(x => {
          groupIds.push(x.ID)
        })
        try {
          await groupDel({ groupIds: groupIds }).then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    async syncPreview() {
      if (this.multipleSelection.length === 0) return
      this.previewLoading = true
      try {
        const groupIds = this.multipleSelection.map(x => x.ID)
        const res = await syncSqlGroupsPreview({ groupIds })
        if (res.data) {
          this.previewResult = res.data
          this.previewDialogVisible = true
        } else {
          Message.error(res.msg || '预览失败')
        }
      } catch (e) {
        Message.error(e.msg || '预览失败')
      } finally {
        this.previewLoading = false
      }
    },
    confirmSyncFromPreview() {
      this.previewDialogVisible = false
      this.batchSync()
    },
    // 批量同步
    batchSync() {
      this.$confirm('此操作批量同步数据到Ldap, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const groupIds = []
        this.multipleSelection.forEach(x => {
          groupIds.push(x.ID)
        })
        try {
          await syncSqlGroups({ groupIds: groupIds }).then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消同步'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      try {
        await groupDel({ groupIds: [Id] }).then(res => {
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },
    // 单个同步
    async singleSync(Id) {
      this.loading = true
      try {
        await syncSqlGroups({ groupIds: [Id] }).then(res => {
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    },
    // treeselect
    normalizer(node) {
      return {
        id: node.ID,
        label: node.groupName,
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
    },
    syncDingTalkDepts() {
      this.loading = true
      syncDingTalkDeptsApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncWeComDepts() {
      this.loading = true
      syncWeComDeptsApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncFeiShuDepts() {
      this.loading = true
      syncFeiShuDeptsApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncOpenLdapDepts() {
      this.loading = true
      syncOpenLdapDeptsApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    }
  }
}
</script>

<style scoped>
  .container-card{
    margin: 10px;
    margin-bottom: 100px;
  }

  .delete-popover{
    margin-left: 10px;
  }
  .sync-preview-body .el-descriptions { margin-bottom: 12px; }
  .sync-preview-body .preview-list { font-size: 13px; margin: 8px 0; word-break: break-all; }
  .sync-preview-body .preview-tip { font-size: 12px; color: #909399; margin-top: 12px; }
   .transfer-footer {
    margin-left: 20px;
    padding: 6px 5px;
  }
</style>
