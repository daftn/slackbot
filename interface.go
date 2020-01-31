package slackbot

import (
	"context"
	"github.com/nlopes/slack"
	"io"
)

// MessagingClient allows for mocking the slack client for testing
type MessagingClient interface {
	AddChannelReminder(string, string, string) (*slack.Reminder, error)
	AddPin(string, slack.ItemRef) error
	AddPinContext(context.Context, string, slack.ItemRef) error
	AddReaction(string, slack.ItemRef) error
	AddReactionContext(context.Context, string, slack.ItemRef) error
	AddStar(string, slack.ItemRef) error
	AddStarContext(context.Context, string, slack.ItemRef) error
	AddUserReminder(string, string, string) (*slack.Reminder, error)
	ArchiveChannel(string) error
	ArchiveChannelContext(context.Context, string) error
	ArchiveConversation(string) error
	ArchiveConversationContext(context.Context, string) error
	ArchiveGroup(string) error
	ArchiveGroupContext(context.Context, string) error
	AuthTest() (*slack.AuthTestResponse, error)
	AuthTestContext(context.Context) (*slack.AuthTestResponse, error)
	CloseConversation(string) (bool, bool, error)
	CloseConversationContext(context.Context, string) (bool, bool, error)
	CloseIMChannel(string) (bool, bool, error)
	CloseIMChannelContext(context.Context, string) (bool, bool, error)
	ConnectRTM() (*slack.Info, string, error)
	ConnectRTMContext(context.Context) (*slack.Info, string, error)
	CreateChannel(string) (*slack.Channel, error)
	CreateChannelContext(context.Context, string) (*slack.Channel, error)
	CreateChildGroup(string) (*slack.Group, error)
	CreateChildGroupContext(context.Context, string) (*slack.Group, error)
	CreateConversation(string, bool) (*slack.Channel, error)
	CreateConversationContext(context.Context, string, bool) (*slack.Channel, error)
	CreateGroup(string) (*slack.Group, error)
	CreateGroupContext(context.Context, string) (*slack.Group, error)
	CreateUserGroup(slack.UserGroup) (slack.UserGroup, error)
	CreateUserGroupContext(context.Context, slack.UserGroup) (slack.UserGroup, error)
	Debug() bool
	Debugf(string, ...interface{})
	Debugln(...interface{})
	DeleteFile(string) error
	DeleteFileComment(string, string) error
	DeleteFileCommentContext(context.Context, string, string) error
	DeleteFileContext(context.Context, string) error
	DeleteMessage(string, string) (string, string, error)
	DeleteMessageContext(context.Context, string, string) (string, string, error)
	DeleteReminder(string) error
	DeleteUserPhoto() error
	DeleteUserPhotoContext(context.Context) error
	DisableUser(string, string) error
	DisableUserContext(context.Context, string, string) error
	DisableUserGroup(string) (slack.UserGroup, error)
	DisableUserGroupContext(context.Context, string) (slack.UserGroup, error)
	Disconnect() error
	EnableUserGroup(string) (slack.UserGroup, error)
	EnableUserGroupContext(context.Context, string) (slack.UserGroup, error)
	EndDND() error
	EndDNDContext(context.Context) error
	EndSnooze() (*slack.DNDStatus, error)
	EndSnoozeContext(context.Context) (*slack.DNDStatus, error)
	GetAccessLogs(slack.AccessLogParameters) ([]slack.Login, *slack.Paging, error)
	GetAccessLogsContext(context.Context, slack.AccessLogParameters) ([]slack.Login, *slack.Paging, error)
	GetBillableInfo(string) (map[string]slack.BillingActive, error)
	GetBillableInfoContext(context.Context, string) (map[string]slack.BillingActive, error)
	GetBillableInfoForTeam() (map[string]slack.BillingActive, error)
	GetBillableInfoForTeamContext(context.Context) (map[string]slack.BillingActive, error)
	GetBotInfo(string) (*slack.Bot, error)
	GetBotInfoContext(context.Context, string) (*slack.Bot, error)
	GetChannel(string) (slack.Channel, error)
	GetChannelHistory(string, slack.HistoryParameters) (*slack.History, error)
	GetChannelHistoryContext(context.Context, string, slack.HistoryParameters) (*slack.History, error)
	GetChannelInfo(string) (*slack.Channel, error)
	GetChannelInfoContext(context.Context, string) (*slack.Channel, error)
	GetChannelReplies(string, string) ([]slack.Message, error)
	GetChannelRepliesContext(context.Context, string, string) ([]slack.Message, error)
	GetChannels(bool, ...slack.GetChannelsOption) ([]slack.Channel, error)
	GetChannelsContext(context.Context, bool, ...slack.GetChannelsOption) ([]slack.Channel, error)
	GetConversationHistory(*slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error)
	GetConversationHistoryContext(context.Context, *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error)
	GetConversationInfo(string, bool) (*slack.Channel, error)
	GetConversationInfoContext(context.Context, string, bool) (*slack.Channel, error)
	GetConversationReplies(*slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error)
	GetConversationRepliesContext(context.Context, *slack.GetConversationRepliesParameters) ([]slack.Message, bool, string, error)
	GetConversations(*slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetConversationsContext(context.Context, *slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetConversationsForUser(*slack.GetConversationsForUserParameters) ([]slack.Channel, string, error)
	GetConversationsForUserContext(context.Context, *slack.GetConversationsForUserParameters) ([]slack.Channel, string, error)
	GetDNDInfo(*string) (*slack.DNDStatus, error)
	GetDNDInfoContext(context.Context, *string) (*slack.DNDStatus, error)
	GetDNDTeamInfo([]string) (map[string]slack.DNDStatus, error)
	GetDNDTeamInfoContext(context.Context, []string) (map[string]slack.DNDStatus, error)
	GetEmoji() (map[string]string, error)
	GetEmojiContext(context.Context) (map[string]string, error)
	GetFile(string, io.Writer) error
	GetFileInfo(string, int, int) (*slack.File, []slack.Comment, *slack.Paging, error)
	GetFileInfoContext(context.Context, string, int, int) (*slack.File, []slack.Comment, *slack.Paging, error)
	GetFiles(slack.GetFilesParameters) ([]slack.File, *slack.Paging, error)
	GetFilesContext(context.Context, slack.GetFilesParameters) ([]slack.File, *slack.Paging, error)
	GetGroupHistory(string, slack.HistoryParameters) (*slack.History, error)
	GetGroupHistoryContext(context.Context, string, slack.HistoryParameters) (*slack.History, error)
	GetGroupInfo(string) (*slack.Group, error)
	GetGroupInfoContext(context.Context, string) (*slack.Group, error)
	GetGroups(bool) ([]slack.Group, error)
	GetGroupsContext(context.Context, bool) ([]slack.Group, error)
	GetIMChannels() ([]slack.IM, error)
	GetIMChannelsContext(context.Context) ([]slack.IM, error)
	GetIMHistory(string, slack.HistoryParameters) (*slack.History, error)
	GetIMHistoryContext(context.Context, string, slack.HistoryParameters) (*slack.History, error)
	GetIncomingEvents() chan slack.RTMEvent
	GetInfo() *slack.Info
	GetPermalink(*slack.PermalinkParameters) (string, error)
	GetPermalinkContext(context.Context, *slack.PermalinkParameters) (string, error)
	GetReactions(slack.ItemRef, slack.GetReactionsParameters) ([]slack.ItemReaction, error)
	GetReactionsContext(context.Context, slack.ItemRef, slack.GetReactionsParameters) ([]slack.ItemReaction, error)
	GetStarred(slack.StarsParameters) ([]slack.StarredItem, *slack.Paging, error)
	GetStarredContext(context.Context, slack.StarsParameters) ([]slack.StarredItem, *slack.Paging, error)
	GetTeamInfo() (*slack.TeamInfo, error)
	GetTeamInfoContext(context.Context) (*slack.TeamInfo, error)
	GetUser(string) (slack.User, error)
	GetUserByEmail(string) (*slack.User, error)
	GetUserByEmailContext(context.Context, string) (*slack.User, error)
	GetUserGroupMembers(string) ([]string, error)
	GetUserGroupMembersContext(context.Context, string) ([]string, error)
	GetUserGroups(...slack.GetUserGroupsOption) ([]slack.UserGroup, error)
	GetUserGroupsContext(context.Context, ...slack.GetUserGroupsOption) ([]slack.UserGroup, error)
	GetUserIdentity() (*slack.UserIdentityResponse, error)
	GetUserIdentityContext(context.Context) (*slack.UserIdentityResponse, error)
	GetUserInfo(string) (*slack.User, error)
	GetUserInfoContext(context.Context, string) (*slack.User, error)
	GetUserPresence(string) (*slack.UserPresence, error)
	GetUserPresenceContext(context.Context, string) (*slack.UserPresence, error)
	GetUserProfile(string, bool) (*slack.UserProfile, error)
	GetUserProfileContext(context.Context, string, bool) (*slack.UserProfile, error)
	GetUsers() ([]slack.User, error)
	GetUsersContext(context.Context) ([]slack.User, error)
	GetUsersInConversation(*slack.GetUsersInConversationParameters) ([]string, string, error)
	GetUsersInConversationContext(context.Context, *slack.GetUsersInConversationParameters) ([]string, string, error)
	GetUsersPaginated(...slack.GetUsersOption) slack.UserPagination
	InviteGuest(string, string, string, string, string) error
	InviteGuestContext(context.Context, string, string, string, string, string) error
	InviteRestricted(string, string, string, string, string) error
	InviteRestrictedContext(context.Context, string, string, string, string, string) error
	InviteToTeam(string, string, string, string) error
	InviteToTeamContext(context.Context, string, string, string, string) error
	InviteUserToChannel(string, string) (*slack.Channel, error)
	InviteUserToChannelContext(context.Context, string, string) (*slack.Channel, error)
	InviteUserToGroup(string, string) (*slack.Group, bool, error)
	InviteUserToGroupContext(context.Context, string, string) (*slack.Group, bool, error)
	InviteUsersToConversation(string, ...string) (*slack.Channel, error)
	InviteUsersToConversationContext(context.Context, string, ...string) (*slack.Channel, error)
	JoinChannel(string) (*slack.Channel, error)
	JoinChannelContext(context.Context, string) (*slack.Channel, error)
	JoinConversation(string) (*slack.Channel, string, []string, error)
	JoinConversationContext(context.Context, string) (*slack.Channel, string, []string, error)
	KickUserFromChannel(string, string) error
	KickUserFromChannelContext(context.Context, string, string) error
	KickUserFromConversation(string, string) error
	KickUserFromConversationContext(context.Context, string, string) error
	KickUserFromGroup(string, string) error
	KickUserFromGroupContext(context.Context, string, string) error
	LeaveChannel(string) (bool, error)
	LeaveChannelContext(context.Context, string) (bool, error)
	LeaveConversation(string) (bool, error)
	LeaveConversationContext(context.Context, string) (bool, error)
	LeaveGroup(string) error
	LeaveGroupContext(context.Context, string) error
	ListFiles(slack.ListFilesParameters) ([]slack.File, *slack.ListFilesParameters, error)
	ListFilesContext(context.Context, slack.ListFilesParameters) ([]slack.File, *slack.ListFilesParameters, error)
	ListPins(string) ([]slack.Item, *slack.Paging, error)
	ListPinsContext(context.Context, string) ([]slack.Item, *slack.Paging, error)
	ListReactions(slack.ListReactionsParameters) ([]slack.ReactedItem, *slack.Paging, error)
	ListReactionsContext(context.Context, slack.ListReactionsParameters) ([]slack.ReactedItem, *slack.Paging, error)
	ListStars(slack.StarsParameters) ([]slack.Item, *slack.Paging, error)
	ListStarsContext(context.Context, slack.StarsParameters) ([]slack.Item, *slack.Paging, error)
	ManageConnection()
	MarkIMChannel(string, string) error
	MarkIMChannelContext(context.Context, string, string) error
	NewOutgoingMessage(string, string, ...slack.RTMsgOption) *slack.OutgoingMessage
	NewRTM(...slack.RTMOption) *slack.RTM
	NewSubscribeUserPresence([]string) *slack.OutgoingMessage
	NewTypingMessage(string) *slack.OutgoingMessage
	OpenConversation(*slack.OpenConversationParameters) (*slack.Channel, bool, bool, error)
	OpenConversationContext(context.Context, *slack.OpenConversationParameters) (*slack.Channel, bool, bool, error)
	OpenDialog(string, slack.Dialog) error
	OpenDialogContext(context.Context, string, slack.Dialog) error
	OpenGroup(string) (bool, bool, error)
	OpenGroupContext(context.Context, string) (bool, bool, error)
	OpenIMChannel(string) (bool, bool, string, error)
	OpenIMChannelContext(context.Context, string) (bool, bool, string, error)
	PostEphemeral(string, string, ...slack.MsgOption) (string, error)
	PostEphemeralContext(context.Context, string, string, ...slack.MsgOption) (string, error)
	PostMessage(string, ...slack.MsgOption) (string, string, error)
	PostMessageContext(context.Context, string, ...slack.MsgOption) (string, string, error)
	RemovePin(string, slack.ItemRef) error
	RemovePinContext(context.Context, string, slack.ItemRef) error
	RemoveReaction(string, slack.ItemRef) error
	RemoveReactionContext(context.Context, string, slack.ItemRef) error
	RemoveStar(string, slack.ItemRef) error
	RemoveStarContext(context.Context, string, slack.ItemRef) error
	RenameChannel(string, string) (*slack.Channel, error)
	RenameChannelContext(context.Context, string, string) (*slack.Channel, error)
	RenameConversation(string, string) (*slack.Channel, error)
	RenameConversationContext(context.Context, string, string) (*slack.Channel, error)
	RenameGroup(string, string) (*slack.Channel, error)
	RenameGroupContext(context.Context, string, string) (*slack.Channel, error)
	RevokeFilePublicURL(string) (*slack.File, error)
	RevokeFilePublicURLContext(context.Context, string) (*slack.File, error)
	Search(string, slack.SearchParameters) (*slack.SearchMessages, *slack.SearchFiles, error)
	SearchContext(context.Context, string, slack.SearchParameters) (*slack.SearchMessages, *slack.SearchFiles, error)
	SearchFiles(string, slack.SearchParameters) (*slack.SearchFiles, error)
	SearchFilesContext(context.Context, string, slack.SearchParameters) (*slack.SearchFiles, error)
	SearchMessages(string, slack.SearchParameters) (*slack.SearchMessages, error)
	SearchMessagesContext(context.Context, string, slack.SearchParameters) (*slack.SearchMessages, error)
	SendAuthRevoke(string) (*slack.AuthRevokeResponse, error)
	SendAuthRevokeContext(context.Context, string) (*slack.AuthRevokeResponse, error)
	SendMessage(*slack.OutgoingMessage)
	SendMessageContext(context.Context, string, ...slack.MsgOption) (string, string, string, error)
	SendSSOBindingEmail(string, string) error
	SendSSOBindingEmailContext(context.Context, string, string) error
	SetChannelPurpose(string, string) (string, error)
	SetChannelPurposeContext(context.Context, string, string) (string, error)
	SetChannelReadMark(string, string) error
	SetChannelReadMarkContext(context.Context, string, string) error
	SetChannelTopic(string, string) (string, error)
	SetChannelTopicContext(context.Context, string, string) (string, error)
	SetGroupPurpose(string, string) (string, error)
	SetGroupPurposeContext(context.Context, string, string) (string, error)
	SetGroupReadMark(string, string) error
	SetGroupReadMarkContext(context.Context, string, string) error
	SetGroupTopic(string, string) (string, error)
	SetGroupTopicContext(context.Context, string, string) (string, error)
	SetPurposeOfConversation(string, string) (*slack.Channel, error)
	SetPurposeOfConversationContext(context.Context, string, string) (*slack.Channel, error)
	SetRegular(string, string) error
	SetRegularContext(context.Context, string, string) error
	SetRestricted(string, string, ...string) error
	SetRestrictedContext(context.Context, string, string, ...string) error
	SetSnooze(int) (*slack.DNDStatus, error)
	SetSnoozeContext(context.Context, int) (*slack.DNDStatus, error)
	SetTopicOfConversation(string, string) (*slack.Channel, error)
	SetTopicOfConversationContext(context.Context, string, string) (*slack.Channel, error)
	SetUltraRestricted(string, string, string) error
	SetUltraRestrictedContext(context.Context, string, string, string) error
	SetUserAsActive() error
	SetUserAsActiveContext(context.Context) error
	SetUserCustomStatus(string, string, int64) error
	SetUserCustomStatusContext(context.Context, string, string, int64) error
	//SetUserCustomStatusContextWithUser(context.Context, string, string, string, int64) error
	//SetUserCustomStatusWithUser(string, string, string, int64) error
	SetUserPhoto(string, slack.UserSetPhotoParams) error
	SetUserPhotoContext(context.Context, string, slack.UserSetPhotoParams) error
	SetUserPresence(string) error
	SetUserPresenceContext(context.Context, string) error
	ShareFilePublicURL(string) (*slack.File, []slack.Comment, *slack.Paging, error)
	ShareFilePublicURLContext(context.Context, string) (*slack.File, []slack.Comment, *slack.Paging, error)
	StartRTM() (*slack.Info, string, error)
	StartRTMContext(context.Context) (*slack.Info, string, error)
	UnArchiveConversation(string) error
	UnArchiveConversationContext(context.Context, string) error
	UnarchiveChannel(string) error
	UnarchiveChannelContext(context.Context, string) error
	UnarchiveGroup(string) error
	UnarchiveGroupContext(context.Context, string) error
	UnfurlMessage(string, string, map[string]slack.Attachment, ...slack.MsgOption) (string, string, string, error)
	UnsetUserCustomStatus() error
	UnsetUserCustomStatusContext(context.Context) error
	UpdateMessage(string, string, ...slack.MsgOption) (string, string, string, error)
	UpdateMessageContext(context.Context, string, string, ...slack.MsgOption) (string, string, string, error)
	UpdateUserGroup(slack.UserGroup) (slack.UserGroup, error)
	UpdateUserGroupContext(context.Context, slack.UserGroup) (slack.UserGroup, error)
	UpdateUserGroupMembers(string, string) (slack.UserGroup, error)
	UpdateUserGroupMembersContext(context.Context, string, string) (slack.UserGroup, error)
	UploadFile(slack.FileUploadParameters) (*slack.File, error)
	UploadFileContext(context.Context, slack.FileUploadParameters) (*slack.File, error)
}
