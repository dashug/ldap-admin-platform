<template>
  <div>
    <page-header title="用户管理" subtitle="管理目录用户、账号状态与第三方同步">
      <template #actions>
        <el-tag type="info" effect="plain" size="small">目录：{{ directoryTypeText }}</el-tag>
        <el-button type="primary" icon="Plus" @click="create">新建用户</el-button>
      </template>
    </page-header>

    <el-card class="container-card" shadow="never">
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input v-model.trim="params.username" prefix-icon="Search" placeholder="搜索用户名" clearable style="width: 200px;" @keyup.enter="search" @clear="search" />
        <el-input v-model.trim="params.nickname" placeholder="昵称" clearable style="width: 130px;" @keyup.enter="search" @clear="search" />
        <el-select v-model.trim="params.status" placeholder="状态" clearable style="width: 110px;" @change="search" @clear="search">
          <el-option label="正常" value="1" />
          <el-option label="禁用" value="2" />
        </el-select>
        <el-select v-model.trim="params.syncState" placeholder="同步状态" clearable style="width: 120px;" @change="search" @clear="search">
          <el-option label="已同步" value="1" />
          <el-option label="未同步" value="2" />
        </el-select>
        <el-button :loading="loading" icon="Search" @click="search">查询</el-button>
        <div class="filter-bar__spacer" />
        <el-button :disabled="multipleSelection.length === 0" :loading="loading" plain type="danger" icon="Delete" @click="batchDelete">批量删除</el-button>
        <el-dropdown trigger="click" @command="handleMoreCommand">
          <el-button plain>更多操作<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
          <template #dropdown><el-dropdown-menu>
            <el-dropdown-item command="batchSync" :disabled="multipleSelection.length === 0" icon="Upload">批量同步</el-dropdown-item>
            <el-dropdown-item command="syncPreview" :disabled="multipleSelection.length === 0" icon="View">同步预览</el-dropdown-item>
            <el-dropdown-item command="batchImport" icon="Upload">批量导入</el-dropdown-item>
            <el-dropdown-item command="export" icon="Download">导出 Excel</el-dropdown-item>
            <el-dropdown-item v-if="syncConfig.ldapEnableSync" divided command="syncLdap" icon="Download">从原 LDAP 同步</el-dropdown-item>
            <el-dropdown-item v-if="syncConfig.dingTalkEnableSync" command="syncDing" icon="Download">从钉钉同步</el-dropdown-item>
            <el-dropdown-item v-if="syncConfig.feiShuEnableSync" command="syncFeishu" icon="Download">从飞书同步</el-dropdown-item>
            <el-dropdown-item v-if="syncConfig.weComEnableSync" command="syncWecom" icon="Download">从企微同步</el-dropdown-item>
          </el-dropdown-menu></template>
        </el-dropdown>
        <el-dropdown trigger="click" @command="handleColumnCommand">
          <el-button plain icon="Operation">列设置</el-button>
          <template #dropdown><el-dropdown-menu class="column-setting-dropdown">
            <el-dropdown-item command="reset"><i class="el-icon-refresh-left" /> 重置为默认</el-dropdown-item>
            <el-dropdown-item divided disabled>显示列（勾选即显示）</el-dropdown-item>
            <el-dropdown-item v-for="col in defaultUserColumns" :key="col.prop" :command="col.prop">
              <el-checkbox :value="columnConfig.visible[col.prop]" @click.prevent>{{ col.label }}</el-checkbox>
            </el-dropdown-item>
          </el-dropdown-menu></template>
        </el-dropdown>
      </div>

      <el-skeleton v-if="loading && tableData.length === 0" :rows="6" animated style="padding: 12px 4px;" />
      <el-table v-else v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange" @header-dragend="handleUserTableHeaderDragend">
        <el-table-column type="selection" width="55" align="center" />
        <template v-for="col in visibleUserColumns">
          <el-table-column
            v-if="col.prop === 'status'"
            :key="col.prop"
            show-overflow-tooltip
            :label="col.label"
            align="center"
            :width="col.width || 70"
            :min-width="col.minWidth"
          >
            <template #default="scope">
              <el-switch v-model="scope.row.status" :active-value="1" :inactive-value="2" @change="userStateChanged(scope.row)" />
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :key="col.prop"
            show-overflow-tooltip
            sortable
            :prop="col.prop"
            :label="col.label"
            :width="col.width"
            :min-width="col.minWidth"
          />
        </template>
        <el-table-column fixed="right" label="操作" align="center" width="190">
          <template #default="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="small" icon="Edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="重置密码" effect="dark" placement="top">
              <el-popconfirm title="确定重置该用户密码吗？" @confirm="openConfirmDialog('resetPassword', { username: scope.row.username }, '确认重置密码')">
                <template #reference><el-button size="small" icon="Key" circle type="warning"  /></template>
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @confirm="openConfirmDialog('singleDelete', { id: scope.row.ID }, '确认删除')">
                <template #reference><el-button size="small" icon="Delete" circle type="danger"  /></template>
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip v-if="scope.row.syncState == 2" class="delete-popover" content="同步" effect="dark" placement="top">
              <el-popconfirm title="确定同步吗？" @confirm="singleSync(scope.row.ID)">
                <template #reference><el-button size="small" icon="Upload" circle type="success"  /></template>
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          :current-page="params.pageNum"
          :page-size="params.pageSize"
          :total="total"
          :page-sizes="[1, 5, 10, 30]"
          layout="total, prev, pager, next, sizes"
          background
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>

      <el-drawer :title="dialogFormTitle" v-model="dialogFormVisible" direction="rtl" size="600px" class="form-drawer">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="80px">
          <el-row>
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input ref="password" v-model.trim="dialogFormData.username" :disabled="disabled" placeholder="用户名（拼音）" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="中文名字" prop="nickname">
                <el-input v-model.trim="dialogFormData.nickname" placeholder="中文名字" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="花名" prop="givenName">
                <el-input v-model.trim="dialogFormData.givenName" placeholder="花名" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="邮箱" prop="mail">
                <el-input v-model.trim="dialogFormData.mail" placeholder="邮箱" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <!-- 修改用户时，不显示密码字段 -->
              <el-form-item v-if="dialogType === 'create'" :label="dialogType === 'create' ? '新密码':'重置密码'" prop="password">
                <el-input v-model.trim="dialogFormData.password" autocomplete="off" :type="passwordType" :placeholder="dialogType === 'create' ? '新密码':'重置密码'" />
                <span class="show-pwd" @click="showPwd">
                  <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
                </span>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="角色" prop="roleIds">
                <el-select v-model.trim="dialogFormData.roleIds" multiple placeholder="请选择角色" style="width:100%">
                  <el-option
                    v-for="item in roles"
                    :key="item.ID"
                    :label="item.name"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="状态" prop="status">
                <el-select v-model.trim="dialogFormData.status" placeholder="请选择状态" style="width:100%">
                  <el-option label="正常" :value="1" />
                  <el-option label="禁用" :value="2" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="手机号" prop="mobile">
                <el-input v-model.trim="dialogFormData.mobile" placeholder="手机号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="工号" prop="jobNumber">
                <el-input v-model.trim="dialogFormData.jobNumber" placeholder="工号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="职位" prop="position">
                <el-input v-model.trim="dialogFormData.position" placeholder="职业" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="所属部门" prop="departmentId">
                <el-tree-select
                  v-model="dialogFormData.departmentId"
                  :data="departmentsOptions"
                  :props="treeSelectProps"
                  node-key="ID"
                  multiple
                  show-checkbox
                  check-strictly
                  :render-after-expand="false"
                  placeholder="请选择部门"
                  style="width:100%"
                  @change="treeselectInput"
                />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="地址" prop="postalAddress">
                <el-input v-model.trim="dialogFormData.postalAddress" type="textarea" placeholder="地址" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="说明" prop="introduction">
                <el-input v-model.trim="dialogFormData.introduction" type="textarea" placeholder="说明" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="过期日" prop="expireAt">
                <el-date-picker v-model="dialogFormData.expireAt" type="date" value-format="yyyy-MM-dd" placeholder="不填则永不过期" clearable style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <template #footer><div class="drawer-footer">
          <el-button @click="cancelForm()">取 消</el-button>
          <el-button :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div></template>
      </el-drawer>

      <!-- 重置密码结果对话框 -->
      <el-dialog
        title="密码重置成功"
        v-model="resetPasswordDialogVisible"
        width="400px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        @close="closeResetPasswordDialog"
      >
        <div style="text-align: center;">
          <el-alert
            title="请保存新密码"
            type="warning"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          />
          <p style="margin-bottom: 10px; font-weight: bold;">用户：{{ resetUsername }}</p>
          <p style="margin-bottom: 20px; color: #606266;">新密码：</p>
          <el-input
            v-model="newPassword"
            readonly
            style="margin-bottom: 20px;"
          >
            <template #append><el-button
              icon="DocumentCopy"
              @click="copyPassword"
            >
              复制
            </el-button></template>
          </el-input>
          <el-alert
            title="请立即保存密码，关闭对话框后将无法再次查看"
            type="info"
            :closable="false"
            show-icon
          />
        </div>
        <template #footer><div class="dialog-footer">
          <el-button type="primary" @click="closeResetPasswordDialog">我已保存</el-button>
        </div></template>
      </el-dialog>

      <!-- 同步预览（Dry Run）结果 -->
      <el-dialog title="同步预览" v-model="previewDialogVisible" width="520px" append-to-body @close="previewResult = null">
        <div v-if="previewResult" class="sync-preview-body">
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="将新增到 LDAP">{{ previewResult.addCount }} 个用户</el-descriptions-item>
            <el-descriptions-item label="LDAP 中已存在（将更新）">{{ previewResult.updateCount }} 个用户</el-descriptions-item>
          </el-descriptions>
          <template v-if="(previewResult.addList && previewResult.addList.length) || (previewResult.updateList && previewResult.updateList.length)">
            <p v-if="previewResult.addList && previewResult.addList.length" class="preview-list">
              <strong>新增：</strong>{{ previewResult.addList.join('、') }}
            </p>
            <p v-if="previewResult.updateList && previewResult.updateList.length" class="preview-list">
              <strong>已存在：</strong>{{ previewResult.updateList.join('、') }}
            </p>
          </template>
          <p class="preview-tip">以上为预览结果，未执行实际同步。点击「执行同步」将正式同步到 LDAP。</p>
        </div>
        <template #footer><div class="dialog-footer">
          <el-button size="small" @click="previewDialogVisible = false">关 闭</el-button>
          <el-button size="small" type="success" :disabled="!previewResult || (previewResult.addCount === 0 && previewResult.updateCount === 0)" @click="confirmSyncFromPreview">执行同步</el-button>
        </div></template>
      </el-dialog>

      <!-- 批量导入用户 -->
      <user-import-dialog ref="importDialog" @done="getTableData" />

      <!-- 敏感操作二次确认：输入确认文案后执行 -->
      <el-dialog title="敏感操作确认" v-model="confirmDialogVisible" width="400px" append-to-body @close="closeConfirmDialog">
        <p class="confirm-dialog-tip">请输入「<strong>{{ confirmPhrase }}</strong>」以继续执行操作。</p>
        <el-input v-model.trim="confirmInput" :placeholder="'请输入 ' + confirmPhrase" clearable @keyup.enter="submitConfirmDialog" />
        <template #footer><div class="dialog-footer">
          <el-button size="small" @click="closeConfirmDialog">取 消</el-button>
          <el-button size="small" type="danger" :disabled="confirmInput !== confirmPhrase" @click="submitConfirmDialog">确 定</el-button>
        </div></template>
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import JSEncrypt from 'jsencrypt'
import { getUsers, createUser, updateUserById, batchDeleteUserByIds, changeUserStatus, syncDingTalkUsersApi, syncWeComUsersApi, syncFeiShuUsersApi, syncOpenLdapUsersApi, syncSqlUsers, syncSqlUsersPreview } from '@/api/personnel/user'
import { resetPassword } from '@/api/system/user'
import { getRoles } from '@/api/system/role'
import { getGroupTree } from '@/api/personnel/group'
import { getConfig } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import UserImportDialog from '@/components/UserImportDialog/index.vue'
import { export_json_to_excel } from '@/vendor/Export2Excel'
import { loadTableColumnConfig, saveTableColumnConfig, defaultUserColumns, STORAGE_KEY_USER_TABLE } from '@/utils/tableColumnSettings'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'User',
  components: {
    PageHeader,
    UserImportDialog
  },
  data() {
    var checkPhone = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('手机号不能为空'))
      } else {
        const reg = /1\d{10}/
        if (reg.test(value)) {
          callback()
        } else {
          return callback(new Error('请输入正确的手机号'))
        }
      }
    }
    return {
      // username 默认不可编辑（新增时可编辑、编辑时禁用）；如需放开请配合后端 ldap-user-name-modify
      disabled: false,
      // 查询参数
      params: {
        username: '',
        nickname: '',
        status: '',
        syncState: '',
        mobile: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,
      exportLoading: false,
      isUpdate: false,
      // 部门信息数据
      treeselectValue: 0,
      // el-tree-select 字段映射（替代旧 vue-treeselect 的 normalizer）
      treeSelectProps: {
        value: 'ID',
        children: 'children',
        label: (data) => data.groupType + '=' + data.groupName,
        disabled: (data) => data.groupType === 'ou' || data.groupName === 'root'
      },
      // 角色
      roles: [],
      // 部门信息
      departmentsOptions: [],

      passwordType: 'password',

      publicKey: import.meta.env.VITE_APP_PUBLIC_KEY,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: '',
        ID: '',
        mail: '',
        givenName: '',
        jobNumber: '',
        postalAddress: '',
        departments: '',
        position: '',
        departmentId: undefined,
        expireAt: ''
      },
      dialogFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        password: [
          { required: false, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 30, message: '长度在 6 到 30 个字符', trigger: 'blur' }
        ],
        mail: [
          { required: true, message: '请输入邮箱', trigger: 'blur' }
        ],
        jobNumber: [
          { required: true, message: '请输入工号', trigger: 'blur' },
          { min: 0, max: 20, message: '长度在 0 到 20 个字符', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '请输入昵称', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        mobile: [
          { required: true, validator: checkPhone, trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ],
        departmentId: [
          { required: true, message: '请选择部门', trigger: 'change' },
          { validator: (rule, value, callBack) => {
            if (value < 1) {
              callBack('请选择有效的部门')
            } else {
              callBack()
            }
          }
          }
        ],
        introduction: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 100 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      changeUserStatusFormData: {
        id: '',
        status: ''
      },

      // 重置密码结果对话框
      resetPasswordDialogVisible: false,
      newPassword: '',
      resetUsername: '',

      // 同步配置（驱动目录类型标签与「更多操作」中的同步入口可见性）
      syncConfig: {
        ldapEnableSync: false,
        dingTalkEnableSync: false,
        feiShuEnableSync: false,
        weComEnableSync: false,
        directoryType: 'openldap'
      },

      // 表格列设置（显示/隐藏、列宽），从 localStorage 恢复
      columnConfig: (() => {
        const c = loadTableColumnConfig(STORAGE_KEY_USER_TABLE, defaultUserColumns)
        return { visible: c.visible, widths: c.widths }
      })(),

      // 敏感操作二次确认弹窗
      confirmDialogVisible: false,
      confirmDialogPayload: { type: '', id: null, username: '' },
      confirmInput: '',
      confirmPhrase: '',

      previewLoading: false,
      previewDialogVisible: false,
      previewResult: null
    }
  },
  created() {
    this.getTableData()
    this.getRoles()
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
    visibleUserColumns() {
      return defaultUserColumns
        .filter(c => this.columnConfig.visible[c.prop] !== false)
        .map(c => ({
          ...c,
          width: this.columnConfig.widths[c.prop] || c.width || undefined,
          minWidth: c.minWidth
        }))
    }
  },
  methods: {
    toggleUserColumnVisible(prop) {
      this.columnConfig.visible[prop] = !this.columnConfig.visible[prop]
      saveTableColumnConfig(STORAGE_KEY_USER_TABLE, this.columnConfig)
    },
    handleColumnCommand(cmd) {
      if (cmd === 'reset') this.resetUserColumnSettings()
      else this.toggleUserColumnVisible(cmd)
    },
    resetUserColumnSettings() {
      this.columnConfig.visible = {}
      this.columnConfig.widths = {}
      defaultUserColumns.forEach(c => {
        this.columnConfig.visible[c.prop] = c.visible !== false
        if (c.width) this.columnConfig.widths[c.prop] = c.width
      })
      saveTableColumnConfig(STORAGE_KEY_USER_TABLE, this.columnConfig)
      Message.success('已重置为默认列')
    },
    handleUserTableHeaderDragend(newWidth, oldWidth, column) {
      if (column && column.property) {
        this.columnConfig.widths[column.property] = newWidth
        saveTableColumnConfig(STORAGE_KEY_USER_TABLE, this.columnConfig)
      }
    },
    openConfirmDialog(type, payload, phrase) {
      this.confirmDialogPayload = { type, ...payload }
      this.confirmPhrase = phrase
      this.confirmInput = ''
      this.confirmDialogVisible = true
    },
    closeConfirmDialog() {
      this.confirmDialogVisible = false
      this.confirmInput = ''
      this.confirmPhrase = ''
      this.confirmDialogPayload = { type: '', id: null, username: '' }
    },
    async submitConfirmDialog() {
      if (this.confirmInput !== this.confirmPhrase) return
      const { type, id, username } = this.confirmDialogPayload
      this.closeConfirmDialog()
      if (type === 'singleDelete' && id) {
        await this.singleDelete(id)
      } else if (type === 'batchDelete') {
        this.loading = true
        try {
          const userIds = this.multipleSelection.map(x => x.ID)
          await batchDeleteUserByIds({ userIds }).then(res => this.judgeResult(res))
          this.getTableData()
        } finally {
          this.loading = false
        }
      } else if (type === 'resetPassword' && username) {
        await this.resetUserPassword(username)
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
    // 「更多操作」下拉
    handleMoreCommand(cmd) {
      if (cmd === 'batchSync') this.batchSync()
      else if (cmd === 'syncPreview') this.syncPreview()
      else if (cmd === 'batchImport') this.$refs.importDialog.open()
      else if (cmd === 'export') this.exportUserList()
      else if (cmd === 'syncLdap') this.syncOpenLdapUsers()
      else if (cmd === 'syncDing') this.syncDingTalkUsers()
      else if (cmd === 'syncFeishu') this.syncFeiShuUsers()
      else if (cmd === 'syncWecom') this.syncWeComUsers()
    },
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getUsers(this.params)
        data.users.forEach(item => {
          const dataStrArr = item.departmentId.split(',')
          const dataIntArr = []
          dataStrArr.forEach(item => {
            dataIntArr.push(+item)
          })
          item.departmentId = dataIntArr
        })
        this.tableData = data.users
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    // 获取所有的分组信息，用于弹框选取上级分组
    async getAllGroups() {
      this.loading = true
      try {
        const checkParams = {
          pageNum: 1,
          pageSize: 1000 // 平常百姓人家应该不会有这么多数据吧
        }
        const { data } = await getGroupTree(checkParams)
        this.departmentsOptions = [{ ID: 0, groupName: '请选择部门信息', groupType: 'T', children: data }]
      } finally {
        this.loading = false
      }
    },
    // 获取角色数据
    async getRoles() {
      const res = await getRoles(null)

      this.roles = res.data.roles
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增用户'
      this.dialogType = 'create'
      this.disabled = false
      this.getAllGroups()
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.disabled = true
      this.getAllGroups()
      this.dialogFormData.ID = row.ID
      this.dialogFormData.username = row.username
      this.dialogFormData.password = ''
      this.dialogFormData.nickname = row.nickname
      this.dialogFormData.status = row.status
      this.dialogFormData.mobile = row.mobile
      this.dialogFormData.introduction = row.introduction
      // 遍历角色数组，获取角色ID
      this.dialogFormData.roleIds = row.roles.map(item => item.ID)

      this.dialogFormTitle = '修改用户'
      this.dialogType = 'update'
      this.passwordType = 'password'
      this.dialogFormVisible = true

      this.dialogFormData.mail = row.mail
      this.dialogFormData.givenName = row.givenName
      this.dialogFormData.jobNumber = row.jobNumber
      this.dialogFormData.postalAddress = row.postalAddress
      this.dialogFormData.departments = row.departments
      this.dialogFormData.departmentId = row.departmentId
      this.dialogFormData.position = row.position
      this.dialogFormData.expireAt = row.expireAt ? (typeof row.expireAt === 'string' ? row.expireAt.slice(0, 10) : row.expireAt) : ''
    },

    // 将 部门id 转换为 部门name
    setDepartmentNameByDepartmentId() {
      const ids = this.dialogFormData.departmentId
      if (!ids || !ids.length) return
      const departments = []
      // 深度优先遍函数
      const dfs = (node, cb) => {
        if (!node) return
        cb(node)
        if (node.children && node.children.length) {
          node.children.forEach(item => {
            dfs(item, cb)
          })
        }
      }
      dfs(this.departmentsOptions[0], node => {
        if (ids.includes(node.ID)) {
          departments.push(node.groupName)
        }
      })
      this.dialogFormData.departments = departments.join(',')
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
      if (this.dialogFormData.nickname === '') {
        Message({
          showClose: true,
          message: '请填写昵称',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.username === '') {
        Message({
          showClose: true,
          message: '请填写用户名',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.mail === '') {
        Message({
          showClose: true,
          message: '请填写邮箱',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.jobNumber === '') {
        Message({
          showClose: true,
          message: '请填写工号',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.mobile === '') {
        Message({
          showClose: true,
          message: '请填写手机号',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.status === '') {
        Message({
          showClose: true,
          message: '请填写状态',
          type: 'error'
        })
        return false
      }
      if (this.dialogFormData.roleIds === '') {
        Message({
          showClose: true,
          message: '请选择角色列表',
          type: 'error'
        })
        return false
      }
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true
          // 在这里自动填充下部门字段
          this.setDepartmentNameByDepartmentId()
          this.dialogFormDataCopy = { ...this.dialogFormData }
          if (this.dialogFormData.password !== '') {
          // 密码RSA加密处理
            const encryptor = new JSEncrypt()
            // 设置公钥
            encryptor.setPublicKey(this.publicKey)
            // 加密密码
            const encPassword = encryptor.encrypt(this.dialogFormData.password)
            this.dialogFormDataCopy.password = encPassword
          }
          try {
            if (this.dialogType === 'create') {
              await createUser(this.dialogFormDataCopy).then(res => {
                this.judgeResult(res)
              })
            } else {
              await updateUserById(this.dialogFormDataCopy).then(res => {
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
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: '',
        departments: '',
        position: '',
        departmentId: undefined
      }
    },

    // 批量删除：先弹确认，再弹二次确认（输入文案）
    batchDelete() {
      this.$confirm('此操作将永久删除选中的用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.openConfirmDialog('batchDelete', {}, '确认批量删除')
      }).catch(() => {
        Message({ showClose: true, type: 'info', message: '已取消删除' })
      })
    },
    async syncPreview() {
      if (this.multipleSelection.length === 0) return
      this.previewLoading = true
      try {
        const userIds = this.multipleSelection.map(x => x.ID)
        const res = await syncSqlUsersPreview({ userIds })
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
      this.$confirm('此操作批量将数据库的用户同步到Ldap, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        try {
          await syncSqlUsers({ userIds: userIds }).then(res => {
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

    /** 导出当前页用户列表为 Excel */
    async exportUserList() {
      this.exportLoading = true
      try {
        const header = ['用户名', '中文名', '花名', '状态', '邮箱', '手机号', '工号', '部门', '职位', '创建人', '说明', 'DN', '创建时间', '更新时间', '最后登录时间', '过期日']
        const data = this.tableData.map(row => [
          row.username || '',
          row.nickname || '',
          row.givenName || '',
          row.status === 1 ? '正常' : row.status === 2 ? '禁用' : String(row.status || ''),
          row.mail || '',
          row.mobile || '',
          row.jobNumber || '',
          row.departments || '',
          row.position || '',
          row.creator || '',
          row.introduction || '',
          row.userDn || '',
          row.CreatedAt || '',
          row.UpdatedAt || '',
          row.lastLoginAt || '',
          row.expireAt ? (typeof row.expireAt === 'string' ? row.expireAt.slice(0, 10) : row.expireAt) : ''
        ])
        const filename = `用户列表_${new Date().toISOString().slice(0, 10).replace(/-/g, '')}.xlsx`
        await export_json_to_excel({ header, data, filename, autoWidth: true, bookType: 'xlsx' })
        Message.success('导出成功')
      } catch (e) {
        Message.error('导出失败：' + (e && e.message ? e.message : String(e)))
      } finally {
        this.exportLoading = false
      }
    },

    // 监听 switch 开关 状态改变
    async userStateChanged(userInfo) {
      this.changeUserStatusFormData.id = userInfo.ID
      this.changeUserStatusFormData.status = userInfo.status
      const { code } = await changeUserStatus(this.changeUserStatusFormData)
      if (code !== 0) {
        return Message.error('更新用户状态失败')
      }
      Message.success('更新用户状态成功')
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      try {
        await batchDeleteUserByIds({ userIds: [Id] }).then(res => {
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
        await syncSqlUsers({ userIds: [Id] }).then(res => {
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
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
    treeselectInput(value) {
      this.treeselectValue = value
    },
    syncDingTalkUsers() {
      this.loading = true
      syncDingTalkUsersApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncWeComUsers() {
      this.loading = true
      syncWeComUsersApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncFeiShuUsers() {
      this.loading = true
      syncFeiShuUsersApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    syncOpenLdapUsers() {
      this.loading = true
      syncOpenLdapUsersApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },

    // 重置用户密码
    async resetUserPassword(username) {
      this.loading = true
      try {
        const res = await resetPassword({ username: username })
        if (res.code === 0) {
          this.newPassword = res.data.newPassword
          this.resetUsername = username
          this.resetPasswordDialogVisible = true
          Message({
            showClose: true,
            message: '密码重置成功',
            type: 'success'
          })
        } else {
          Message({
            showClose: true,
            message: res.msg || '密码重置失败',
            type: 'error'
          })
        }
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    // 复制密码到剪贴板
    copyPassword() {
      const textArea = document.createElement('textarea')
      textArea.value = this.newPassword
      document.body.appendChild(textArea)
      textArea.select()
      try {
        document.execCommand('copy')
        Message({
          showClose: true,
          message: '密码已复制到剪贴板',
          type: 'success'
        })
      } catch (err) {
        Message({
          showClose: true,
          message: '复制失败，请手动复制',
          type: 'error'
        })
      }
      document.body.removeChild(textArea)
    },

    // 关闭重置密码对话框
    closeResetPasswordDialog() {
      this.resetPasswordDialogVisible = false
      this.newPassword = ''
      this.resetUsername = ''
    }
  }
}
</script>

<style scoped>
  .container-card {
    margin: 10px;
    margin-bottom: 100px;
  }

  .confirm-dialog-tip {
    margin-bottom: 12px;
    color: #606266;
  }

  .sync-preview-body .el-descriptions { margin-bottom: 12px; }
  .sync-preview-body .preview-list { font-size: 13px; margin: 8px 0; word-break: break-all; }
  .sync-preview-body .preview-tip { font-size: 12px; color: #909399; margin-top: 12px; }

  .toolbar-section {
    margin-bottom: 12px;
  }
  .toolbar-section:last-of-type {
    margin-bottom: 16px;
  }
  .toolbar-section--secondary {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px 12px;
    padding: 10px 12px;
    background: #fafafa;
    border-radius: 4px;
    border: 1px solid #ebeef5;
  }
  .toolbar-form {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
  }
  .toolbar-label {
    font-size: 12px;
    color: #909399;
    margin-right: 4px;
  }
  .toolbar-label--sync {
    margin-left: 8px;
    padding-left: 8px;
    border-left: 1px solid #dcdfe6;
  }
  .toolbar-tag {
    margin-left: 4px;
  }

  .delete-popover {
    margin-left: 10px;
  }

  .pagination-wrap {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    margin-top: 12px;
    margin-bottom: 12px;
  }

  .dialog-footer {
    text-align: right;
  }
  .dialog-footer .el-button {
    margin-left: 10px;
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 3px;
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
  }
</style>
